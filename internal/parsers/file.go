package parsers

import (
	"regexp"
	"strings"

	"github.com/jromero/ugo/internal/tasks"
	"github.com/jromero/ugo/internal/types"
)

var taskFileToken = regexp.MustCompile(`^file=([^;]+);?$`)

type FileParser struct{}

func (f *FileParser) AttemptParse(taskDefinition, nextCodeBlock string) (types.Task, error) {
	if fileMatch := taskFileToken.FindStringSubmatch(taskDefinition); len(fileMatch) > 0 {
		return tasks.NewFileTask(strings.TrimSpace(fileMatch[1]), nextCodeBlock), nil
	}

	return nil, nil
}
