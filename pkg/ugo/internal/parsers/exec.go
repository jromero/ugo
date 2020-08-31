package parsers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var taskExecToken = regexp.MustCompile(`^exec;?(exit-code=(-?[0-9]+))?;?$`)

var _ types.Parser = (*ExecParser)(nil)

type ExecParser struct{}

func (f *ExecParser) AttemptParse(scope, taskDefinition, nextCodeBlock string) (types.Task, error) {
	if execMatch := taskExecToken.FindStringSubmatch(taskDefinition); len(execMatch) > 0 {
		exitCode := 0
		if execMatch[2] != "" {
			var err error
			exitCode, err = strconv.Atoi(execMatch[2])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("parsing weight: %s: %s", execMatch[2], err.Error()))
			}
		}

		return tasks.NewExecTask(scope, nextCodeBlock, exitCode), nil
	}

	return nil, nil
}
