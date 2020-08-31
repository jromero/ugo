package parsers

import (
	"regexp"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var taskAssertContainsToken = regexp.MustCompile(`^assert=contains;?$`)

var _ types.Parser = (*AssertContainsParser)(nil)

type AssertContainsParser struct{}

func (a *AssertContainsParser) AttemptParse(scope, taskDefinition, nextCodeBlock string) (types.Task, error) {
	if taskAssertContainsToken.MatchString(taskDefinition) {
		return tasks.NewAssertContainsTask(scope, nextCodeBlock), nil
	}

	return nil, nil
}
