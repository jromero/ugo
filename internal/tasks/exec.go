package tasks

import (
	"github.com/jromero/ugo/internal/types"
)

var _ types.Task = (*ExecTask)(nil)

type ExecTask struct {
	contents string
	exitCode int
}

func NewExecTask(contents string, exitCode int) *ExecTask {
	return &ExecTask{contents: contents, exitCode: exitCode}
}

func (e *ExecTask) Name() string {
	return "exec"
}

func (e *ExecTask) Contents() string {
	return e.contents
}

func (e *ExecTask) ExitCode() int {
	return e.exitCode
}
