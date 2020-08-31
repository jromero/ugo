package ugo

import "github.com/jromero/ugo/internal"

type Plan struct {
	suites []Suite
}

func NewPlan(suites []Suite) *Plan {
	return &Plan{suites: suites}
}

type Suite struct {
	name   string
	weight int
	tasks  []internal.Task
}

func NewSuite(name string, weight int, tasks []internal.Task) *Suite {
	return &Suite{name: name, weight: weight, tasks: tasks}
}
