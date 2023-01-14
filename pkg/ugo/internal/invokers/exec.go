package invokers

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Invoker = (*ExecInvoker)(nil)

type ExecInvoker struct {
	Logger types.Logger
}

func NewExecInvoker(logger types.Logger) *ExecInvoker {
	return &ExecInvoker{Logger: logger}
}

func (e *ExecInvoker) Supports(task types.Task) bool {
	_, ok := task.(*tasks.ExecTask)
	return ok
}

func (e *ExecInvoker) Invoke(task types.Task, keep bool, workDir, _ string) (output string, err error) {
	return e.executeExec(keep, workDir, task.(*tasks.ExecTask))
}

func (e *ExecInvoker) executeExec(keep bool, workDir string, task *tasks.ExecTask) (output string, err error) {
	tmpScript := filepath.Join(workDir, fmt.Sprintf(".script-%x", sha256.Sum256([]byte(task.Contents()))))
	if !keep {
		defer os.Remove(tmpScript)
	}

	exitCode := task.ExitCode()

	err = ioutil.WriteFile(tmpScript, []byte(task.Contents()), os.ModePerm)
	if err != nil {
		return output, err
	}

	// override umask
	err = os.Chmod(tmpScript, 0755)
	if err != nil {
		return output, err
	}

	bashExe, err := exec.LookPath("bash")
	if err != nil {
		return output, err
	}

	var outBuf bytes.Buffer

	cmd := exec.Cmd{
		Dir:    workDir,
		Path:   bashExe,
		Args:   []string{bashExe, "-e", tmpScript},
		Stdout: &outBuf,
		Stderr: &outBuf,
	}

	e.Logger.Debug("Executing the following:\n%s", task.Contents())
	err = cmd.Run()

	output = outBuf.String()
	e.Logger.Debug("Output:\n%s", output)
	if exitError, ok := err.(*exec.ExitError); err != nil && ok {
		if exitCode == -1 || exitError.ExitCode() == exitCode {
			return output, nil
		}
	}

	return output, err
}
