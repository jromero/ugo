package tasks

import (
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Task = (*FileTask)(nil)

type FileTask struct {
	filename string
	contents string
}

func NewFileTask(filename string, contents string) *FileTask {
	return &FileTask{filename: filename, contents: contents}
}

func (f *FileTask) Name() string {
	return "file"
}

func (f *FileTask) Filename() string {
	return f.filename
}

func (f *FileTask) Contents() string {
	return f.contents
}
