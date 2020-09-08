package tasks

import (
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Task = (*AssertContainsTask)(nil)

type AssertContainsTask struct {
	scope       string
	expected    string
	ignoreLines string
}

func NewAssertContainsTask(scope, expected string, opts ...func(task *AssertContainsTask)) *AssertContainsTask {
	task := &AssertContainsTask{expected: expected, scope: scope}
	for _, opt := range opts {
		opt(task)
	}

	return task
}

func (a *AssertContainsTask) Name() string {
	return "assert:contains"
}

func (a *AssertContainsTask) Scope() string {
	return a.scope
}

func (a *AssertContainsTask) Expected() string {
	return a.expected
}

func (a *AssertContainsTask) IgnoreLines() string {
	return a.ignoreLines
}

func WithIgnoreLines(chars string) func(task *AssertContainsTask) {
	return func(task *AssertContainsTask) {
		task.ignoreLines = chars
	}
}
