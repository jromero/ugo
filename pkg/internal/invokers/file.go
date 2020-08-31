package invokers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jromero/ugo/pkg/internal/tasks"
	"github.com/jromero/ugo/pkg/types"
)

var _ types.Invoker = (*FileInvoker)(nil)

type FileInvoker struct{}

func (f *FileInvoker) Supports(task types.Task) bool {
	_, ok := task.(*tasks.FileTask)
	return ok
}

func (f *FileInvoker) Invoke(task types.Task, workDir, _ string) (output string, err error) {
	return "", executeFile(workDir, task.(*tasks.FileTask))
}

func executeFile(workDir string, task *tasks.FileTask) error {
	if task.Filename() == "" {
		return errors.New("filename for a file task must be provided")
	}

	log.Printf("Writing file (%s) with contents:\n%s", task.Filename(), task.Contents())
	return ioutil.WriteFile(filepath.Join(workDir, fmt.Sprintf("%v", task.Filename())), []byte(task.Contents()), os.ModePerm)
}
