package tasks

import "github.com/jromero/ugo/internal"

var _ internal.Task = (*AssertContainsTask)(nil)

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
