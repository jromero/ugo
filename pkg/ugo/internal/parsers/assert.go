package parsers

import (
	"regexp"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var taskAssertContainsToken = regexp.MustCompile(`^assert=contains;?(ignore-lines=([^;]+))?;?$`)

var _ types.Parser = (*AssertContainsParser)(nil)

type AssertContainsParser struct{}

func (a *AssertContainsParser) AttemptParse(scope, taskDefinition, nextCodeBlock string) (types.Task, error) {
	if execMatch := taskAssertContainsToken.FindStringSubmatch(taskDefinition); len(execMatch) > 0 {
		ignoreLines := ""
		if execMatch[2] != "" {
			ignoreLines = execMatch[2]
		}

		return tasks.NewAssertContainsTask(scope, nextCodeBlock, tasks.WithIgnoreLines(ignoreLines)), nil
	}

	return nil, nil
}
