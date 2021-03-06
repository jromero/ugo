package ugo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jromero/ugo/pkg/ugo/internal/parsers"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var (
	suiteToken        = regexp.MustCompile(`<!--\s*test:suite=([^;]+?);?(weight=([0-9]+))?;?\s*-->`)
	taskPrefixToken   = regexp.MustCompile(`<!--\s*test:((setup|teardown):)?(.*)\s*-->`)
	taskContentsToken = regexp.MustCompile(`(?s)\x60\x60\x60.*?\n(.+?)\n\x60\x60\x60`)
	taskParsers       = []types.Parser{
		&parsers.ExecParser{},
		&parsers.FileParser{},
		&parsers.AssertContainsParser{},
	}
)

var NoSuiteError = errors.New("no suite found")

func Parse(content string) (types.Plan, error) {
	var suites []types.Suite

	suiteSubmatches := suiteToken.FindAllStringSubmatchIndex(content, -1)
	if len(suiteSubmatches) == 0 {
		return types.Plan{}, NoSuiteError
	}

	for i, suiteSubmatch := range suiteSubmatches {
		var (
			name   = content[suiteSubmatch[2]:suiteSubmatch[3]]
			tasks  []types.Task
			weight int
			err    error
		)

		if suiteSubmatch[6] != -1 && suiteSubmatch[7] != -1 {
			value := content[suiteSubmatch[6]:suiteSubmatch[7]]
			weight, err = strconv.Atoi(value)
			if err != nil {
				return types.Plan{}, errors.New(fmt.Sprintf("parsing weight: %s: %s", value, err.Error()))
			}
		}

		// if there is a next suite only search within current section
		var additionalTasks []types.Task
		if i+1 < len(suiteSubmatches) {
			additionalTasks, err = parseTasks(content[suiteSubmatch[1]:suiteSubmatches[i+1][0]])
		} else {
			additionalTasks, err = parseTasks(content[suiteSubmatch[1]:])
		}
		if err != nil {
			return types.Plan{}, err
		}

		suites = append(suites, types.NewSuite(name, weight, append(tasks, additionalTasks...)))
	}

	return types.NewPlan(aggregateSuites(suites)), nil
}

func parseTasks(content string) (tasks []types.Task, err error) {
	if taskSubmatches := taskPrefixToken.FindAllStringSubmatchIndex(content, -1); len(taskSubmatches) > 0 {
		for i, taskSubmatch := range taskSubmatches {
			// if there is a next task only search within current section
			var task types.Task
			if i+1 < len(taskSubmatches) {
				task, err = parseTask(content[taskSubmatch[0]:taskSubmatches[i+1][0]])
				if err != nil {
					return nil, err
				}
			} else {
				task, err = parseTask(content[taskSubmatch[0]:])
				if err != nil {
					return nil, err
				}
			}
			if task != nil {
				tasks = append(tasks, task)
			}
		}
	}

	return tasks, nil
}

func parseTask(content string) (types.Task, error) {
	taskSubmatch := taskPrefixToken.FindStringSubmatchIndex(content)
	if len(taskSubmatch) == 0 {
		return nil, nil
	}

	scope := types.ScopeDefault
	if taskSubmatch[4] != -1 && taskSubmatch[5] != -1 {
		scope = content[taskSubmatch[4]:taskSubmatch[5]]
	}

	taskDefinition := strings.TrimSpace(content[taskSubmatch[6]:taskSubmatch[7]])
	contentAfterTask := content[taskSubmatch[1]:]

	for _, taskParser := range taskParsers {
		task, err := taskParser.AttemptParse(scope, taskDefinition, parseCodeBlock(contentAfterTask))
		if err != nil {
			return nil, err
		}

		if task != nil {
			return task, nil
		}
	}

	return nil, fmt.Errorf("unknown task '%s'", taskDefinition)
}

// parseCodeBlock parses the first instance of a code block and returns it's content
func parseCodeBlock(content string) string {
	taskContentsSubmatch := taskContentsToken.FindStringSubmatch(content)
	if len(taskContentsSubmatch) == 0 {
		return ""
	}

	return taskContentsSubmatch[1]
}
