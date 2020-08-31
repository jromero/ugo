package parsers

import (
	"regexp"
	"strings"

	"github.com/jromero/ugo/internal"
	"github.com/jromero/ugo/internal/tasks"
)

var taskFileToken = regexp.MustCompile(`^file=([^;]+);?$`)

type FileParser struct{}

func (f *FileParser) AttemptParse(taskDefinition, nextCodeBlock string) (internal.Task, error) {
	if fileMatch := taskFileToken.FindStringSubmatch(taskDefinition); len(fileMatch) > 0 {
		return tasks.NewFileTask(strings.TrimSpace(fileMatch[1]), nextCodeBlock), nil
	}

	return nil, nil
}
