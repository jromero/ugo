package parsers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/jromero/ugo/internal"
	"github.com/jromero/ugo/internal/tasks"
)

var taskExecToken = regexp.MustCompile(`^exec;?(exit-code=(-?[0-9]+))?;?$`)

type ExecParser struct{}

func (f *ExecParser) AttemptParse(taskDefinition, nextCodeBlock string) (internal.Task, error) {
	if execMatch := taskExecToken.FindStringSubmatch(taskDefinition); len(execMatch) > 0 {
		exitCode := 0
		if execMatch[2] != "" {
			var err error
			exitCode, err = strconv.Atoi(execMatch[2])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("parsing weight: %s: %s", execMatch[2], err.Error()))
			}
		}

		return tasks.NewExecTask(nextCodeBlock, exitCode), nil
	}

	return nil, nil
}
