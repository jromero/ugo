package invokers

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Invoker = (*ExecInvoker)(nil)

type ExecInvoker struct{}

func (e *ExecInvoker) Supports(task types.Task) bool {
	_, ok := task.(*tasks.ExecTask)
	return ok
}

func (e *ExecInvoker) Invoke(task types.Task, workDir, _ string) (output string, err error) {
	return executeExec(workDir, task.(*tasks.ExecTask))
}

func executeExec(workDir string, task *tasks.ExecTask) (output string, err error) {
	tmpScript := filepath.Join(workDir, fmt.Sprintf(".script-%x", sha256.Sum256([]byte(task.Contents()))))
	defer os.Remove(tmpScript)

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

	log.Printf("Executing the following:\n%s", task.Contents())
	err = cmd.Run()

	output = outBuf.String()
	log.Printf("Output:\n%s", output)
	if exitError, ok := err.(*exec.ExitError); err != nil && ok {
		if exitCode == -1 || exitError.ExitCode() == exitCode {
			return output, nil
		}
	}

	return output, err
}
