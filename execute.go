package ugo

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jromero/ugo/internal"
	"github.com/jromero/ugo/internal/tasks"
)

func Execute(plan Plan) error {
	prevFlags := log.Flags()
	defer func() { log.SetFlags(prevFlags) }()
	prevPrefix := log.Prefix()
	defer func() { log.SetPrefix(prevPrefix) }()

	log.SetFlags(0)

	for _, suite := range plan.suites {
		log.SetPrefix(fmt.Sprintf("[%s] ", suite.name))
		log.Printf("Suite '%s' executing...", suite.name)

		workDir, err := ioutil.TempDir("", fmt.Sprintf("suite-%s-*", suite.name))
		if err != nil {
			return err
		}
		log.Println("Working directory:", workDir)

		var aggrOutput string
		for i, task := range suite.tasks {
			log.SetPrefix(fmt.Sprintf("[%s][task#%d] ", suite.name, i+1))
			log.Printf("--> Running task #%d", i+1)
			output, err := executeTask(aggrOutput, workDir, i+1, task)
			if err != nil {
				return err
			}

			aggrOutput += output
		}
	}

	return nil
}

var ansiPattern = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func executeTask(priorOutput, workDir string, index int, task internal.Task) (output string, err error) {
	switch task.(type) {
	case *tasks.FileTask:
		err = executeFile(workDir, task.(*tasks.FileTask))
	case *tasks.ExecTask:
		output, err = executeExec(workDir, task.(*tasks.ExecTask))
	case *tasks.AssertContainsTask:
		err = executeAssertContains(priorOutput, task.(*tasks.AssertContainsTask))
	default:
		err = errors.New("unknown task: " + task.Name())
	}

	if err != nil {
		return output, fmt.Errorf("task #%d (%s) failed: %s", index, task.Name(), err)
	}

	return output, nil
}

func executeAssertContains(priorOutput string, task *tasks.AssertContainsTask) error {
	expected := task.Expected()
	log.Printf("Checking that output contained:\n%s", expected)
	if !strings.Contains(ansiPattern.ReplaceAllString(priorOutput, ""), expected) {
		return fmt.Errorf("no output contained:\n%s", expected)
	}

	return nil
}

func executeFile(workDir string, task *tasks.FileTask) error {
	if task.Filename() == "" {
		return errors.New("filename for a file task must be provided")
	}

	log.Printf("Writing file (%s) with contents:\n%s", task.Filename(), task.Contents())
	return ioutil.WriteFile(filepath.Join(workDir, fmt.Sprintf("%v", task.Filename())), []byte(task.Contents()), os.ModePerm)
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
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitCode == -1 || exitError.ExitCode() == exitCode {
			return output, nil
		}
	}

	return output, err
}
