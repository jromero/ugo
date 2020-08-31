package tasks

import (
	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Task = (*AssertContainsTask)(nil)

type AssertContainsTask struct {
	expected string
	scope    string
}

func NewAssertContainsTask(scope, expected string) *AssertContainsTask {
	return &AssertContainsTask{expected: expected, scope: scope}
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
