package invokers

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jromero/ugo/internal/tasks"
	"github.com/jromero/ugo/internal/types"
)

var ansiPattern = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

var _ types.Invoker = (*AssertInvoker)(nil)

type AssertInvoker struct{}

func (a *AssertInvoker) Supports(task types.Task) bool {
	_, ok := task.(*tasks.AssertContainsTask)
	return ok
}

func (a *AssertInvoker) Invoke(task types.Task, _, prevOutput string) (output string, err error) {
	containsTask := task.(*tasks.AssertContainsTask)
	return "", executeAssertContains(prevOutput, containsTask)
}

func executeAssertContains(priorOutput string, task *tasks.AssertContainsTask) error {
	expected := task.Expected()
	log.Printf("Checking that output contained:\n%s", expected)
	if !strings.Contains(ansiPattern.ReplaceAllString(priorOutput, ""), expected) {
		return fmt.Errorf("no output contained:\n%s", expected)
	}

	return nil
}
