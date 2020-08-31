package parsers

import (
	"regexp"
	"strings"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var taskFileToken = regexp.MustCompile(`^file=([^;]+);?$`)

var _ types.Parser = (*FileParser)(nil)

type FileParser struct{}

func (f *FileParser) AttemptParse(scope, taskDefinition, nextCodeBlock string) (types.Task, error) {
	if fileMatch := taskFileToken.FindStringSubmatch(taskDefinition); len(fileMatch) > 0 {
		return tasks.NewFileTask(scope, strings.TrimSpace(fileMatch[1]), nextCodeBlock), nil
	}

	return nil, nil
}
