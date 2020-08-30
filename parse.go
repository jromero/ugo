package ugo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	suiteToken              = regexp.MustCompile(`<!--\s*test:suite=([^;]+?);?(weight=([0-9]+))?;?\s*-->`)
	taskPrefixToken         = regexp.MustCompile(`<!--\s*test:(.*)\s*-->`)
	taskFileToken           = regexp.MustCompile(`^file=([^;]+);?$`)
	taskExecToken           = regexp.MustCompile(`^exec;?(exit-code=(-?[0-9]+))?;?$`)
	taskAssertContainsToken = regexp.MustCompile(`^assert=contains;?$`)
	taskContentsToken       = regexp.MustCompile(`(?s)\x60\x60\x60.*?\n(.+?)\n\x60\x60\x60`)
)

var NoSuiteError = errors.New("no suite found")

func Parse(content string) (Plan, error) {
	var suites []Suite

	suiteSubmatches := suiteToken.FindAllStringSubmatchIndex(content, -1)
	if len(suiteSubmatches) == 0 {
		return Plan{}, NoSuiteError
	}

	for i, suiteSubmatch := range suiteSubmatches {
		var (
			name   = content[suiteSubmatch[2]:suiteSubmatch[3]]
			tasks  []Task
			weight int
			err    error
		)

		if suiteSubmatch[6] != -1 && suiteSubmatch[7] != -1 {
			value := content[suiteSubmatch[6]:suiteSubmatch[7]]
			weight, err = strconv.Atoi(value)
			if err != nil {
				return Plan{}, errors.New(fmt.Sprintf("parsing weight: %s: %s", value, err.Error()))
			}
		}

		// if there is a next suite only search within current section
		var additionalTasks []Task
		if i+1 < len(suiteSubmatches) {
			additionalTasks, err = parseTasks(content[suiteSubmatch[1]:suiteSubmatches[i+1][0]])
		} else {
			additionalTasks, err = parseTasks(content[suiteSubmatch[1]:])
		}
		if err != nil {
			return Plan{}, err
		}

		suites = append(suites, *NewSuite(name, weight, append(tasks, additionalTasks...)))
	}

	return *NewPlan(aggregateSuites(suites)), nil
}

func parseTasks(content string) (tasks []Task, err error) {
	if taskSubmatches := taskPrefixToken.FindAllStringSubmatchIndex(content, -1); len(taskSubmatches) > 0 {
		for i, taskSubmatch := range taskSubmatches {
			// if there is a next task only search within current section
			var task *Task
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
				tasks = append(tasks, *task)
			}
		}
	}

	return tasks, nil
}

// parseCodeBlock parses the first instance of a code block and returns it's content
func parseCodeBlock(content string) string {
	taskContentsSubmatch := taskContentsToken.FindStringSubmatch(content)
	if len(taskContentsSubmatch) == 0 {
		return ""
	}

	return taskContentsSubmatch[1]
}

func parseTask(content string) (*Task, error) {
	taskSubmatch := taskPrefixToken.FindStringSubmatchIndex(content)
	if len(taskSubmatch) == 0 {
		return nil, nil
	}

	taskDefinition := strings.TrimSpace(content[taskSubmatch[2]:taskSubmatch[3]])

	// determine type
	if fileMatch := taskFileToken.FindStringSubmatch(taskDefinition); len(fileMatch) > 0 {
		return NewFileTask(strings.TrimSpace(fileMatch[1]), parseCodeBlock(content[taskSubmatch[1]:])), nil
	} else if execMatch := taskExecToken.FindStringSubmatch(taskDefinition); len(execMatch) > 0 {
		exitCode := 0
		if execMatch[2] != "" {
			var err error
			exitCode, err = strconv.Atoi(execMatch[2])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("parsing weight: %s: %s", execMatch[2], err.Error()))
			}
		}

		return NewExecTask(parseCodeBlock(content[taskSubmatch[1]:]), exitCode), nil
	} else if taskAssertContainsToken.MatchString(taskDefinition) {
		return NewAssertContainsTask(parseCodeBlock(content[taskSubmatch[1]:])), nil
	}

	return nil, fmt.Errorf("unknown task '%s'", taskDefinition)
}
