package tasks

import (
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Task = (*FileTask)(nil)

type FileTask struct {
	filename string
	contents string
	scope    string
}

func NewFileTask(scope string, filename string, contents string) *FileTask {
	return &FileTask{filename: filename, contents: contents, scope: scope}
}

func (f *FileTask) Name() string {
	return "file"
}

func (f *FileTask) Scope() string {
	return f.scope
}

func (f *FileTask) Filename() string {
	return f.filename
}

func (f *FileTask) Contents() string {
	return f.contents
}
