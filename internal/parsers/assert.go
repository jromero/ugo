package parsers

import (
	"regexp"

	"github.com/jromero/ugo/internal"
	"github.com/jromero/ugo/internal/tasks"
)

var taskAssertContainsToken = regexp.MustCompile(`^assert=contains;?$`)

type AssertContainsParser struct{}

func (a *AssertContainsParser) AttemptParse(taskDefinition, nextCodeBlock string) (internal.Task, error) {
	if taskAssertContainsToken.MatchString(taskDefinition) {
		return tasks.NewAssertContainsTask(nextCodeBlock), nil
	}

	return nil, nil
}
