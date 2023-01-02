package invokers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var ansiPattern = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

var _ types.Invoker = (*AssertInvoker)(nil)

type AssertInvoker struct {
	Logger types.Logger
}

func NewAssertInvoker(logger types.Logger) *AssertInvoker {
	return &AssertInvoker{Logger: logger}
}

func (a *AssertInvoker) Supports(task types.Task) bool {
	_, ok := task.(*tasks.AssertContainsTask)
	return ok
}

func (a *AssertInvoker) Invoke(task types.Task, keep bool, _, prevOutput string) (output string, err error) {
	containsTask := task.(*tasks.AssertContainsTask)
	return "", a.executeAssertContains(prevOutput, containsTask)
}

func (a *AssertInvoker) executeAssertContains(priorOutput string, task *tasks.AssertContainsTask) error {
	expected := task.Expected()

	sanitized := ansiPattern.ReplaceAllString(priorOutput, "")
	if task.IgnoreLines() != "" {
		rx := regexp.MustCompile(fmt.Sprintf(`(?m)^%s(\n|$)`, regexp.QuoteMeta(task.IgnoreLines())))
		expected = regexp.QuoteMeta(rx.ReplaceAllString(expected, "~~WILDCARD~~${1}"))
		expected = strings.ReplaceAll(expected, "~~WILDCARD~~", ".*")
		expected = "(?ms).*" + expected + ".*"

		a.Logger.Debug("Checking that output matches pattern:\n%s", expected)
		if !regexp.MustCompile(expected).MatchString(sanitized) {
			a.Logger.Info("Content:\n%s", sanitized)
			return fmt.Errorf("Output didn't match pattern:\n%s", expected)
		}
	} else {
		a.Logger.Debug("Checking that output contained:\n%s", expected)

		if !strings.Contains(sanitized, expected) {
			dmp := diffmatchpatch.New()
			diff := dmp.DiffPrettyText(dmp.DiffMain(sanitized, expected, false))
			a.Logger.Info("Diff between output and expected:\n%s", diff)
			return fmt.Errorf("Output did not contain:\n%s", expected)
		}
	}

	return nil
}
