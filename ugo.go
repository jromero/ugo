package ugo

import (
	"github.com/jromero/ugo/internal/types"
)

type Plan struct {
	suites []Suite
}

func NewPlan(suites []Suite) *Plan {
	return &Plan{suites: suites}
}

type Suite struct {
	name   string
	weight int
	tasks  []types.Task
}

func NewSuite(name string, weight int, tasks []types.Task) *Suite {
	return &Suite{name: name, weight: weight, tasks: tasks}
}
