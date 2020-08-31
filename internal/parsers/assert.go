package parsers

import (
	"regexp"

	"github.com/jromero/ugo/internal/tasks"
	"github.com/jromero/ugo/internal/types"
)

var taskAssertContainsToken = regexp.MustCompile(`^assert=contains;?$`)

type AssertContainsParser struct{}

func (a *AssertContainsParser) AttemptParse(taskDefinition, nextCodeBlock string) (types.Task, error) {
	if taskAssertContainsToken.MatchString(taskDefinition) {
		return tasks.NewAssertContainsTask(nextCodeBlock), nil
	}

	return nil, nil
}
