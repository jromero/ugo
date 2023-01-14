package invokers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Invoker = (*FileInvoker)(nil)

type FileInvoker struct {
	Logger types.Logger
}

func NewFileInvoker(logger types.Logger) *FileInvoker {
	return &FileInvoker{
		Logger: logger,
	}
}

func (f *FileInvoker) Supports(task types.Task) bool {
	_, ok := task.(*tasks.FileTask)
	return ok
}

func (f *FileInvoker) Invoke(task types.Task, keep bool, workDir, _ string) (output string, err error) {
	return "", f.executeFile(workDir, task.(*tasks.FileTask))
}

func (f *FileInvoker) executeFile(workDir string, task *tasks.FileTask) error {
	if task.Filename() == "" {
		return errors.New("filename for a file task must be provided")
	}

	f.Logger.Debug("Writing file (%s) with contents:\n%s", task.Filename(), task.Contents())
	return ioutil.WriteFile(filepath.Join(workDir, fmt.Sprintf("%v", task.Filename())), []byte(task.Contents()), os.ModePerm)
}
