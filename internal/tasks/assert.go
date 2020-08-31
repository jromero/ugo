package tasks

import (
	"github.com/jromero/ugo/internal/types"
)

var _ types.Task = (*AssertContainsTask)(nil)

type AssertContainsTask struct {
	expected string
}

func NewAssertContainsTask(expected string) *AssertContainsTask {
	return &AssertContainsTask{expected: expected}
}

func (a *AssertContainsTask) Name() string {
	return "assert:contains"
}

func (a *AssertContainsTask) Expected() string {
	return a.expected
}
