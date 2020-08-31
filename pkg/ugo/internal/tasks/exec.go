package tasks

import (
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Task = (*ExecTask)(nil)

type ExecTask struct {
	contents string
	exitCode int
	scope    string
}

func NewExecTask(scope, contents string, exitCode int) *ExecTask {
	return &ExecTask{contents: contents, exitCode: exitCode, scope: scope}
}

func (e *ExecTask) Name() string {
	return "exec"
}

func (e *ExecTask) Scope() string {
	return e.scope
}

func (e *ExecTask) Contents() string {
	return e.contents
}

func (e *ExecTask) ExitCode() int {
	return e.exitCode
}
