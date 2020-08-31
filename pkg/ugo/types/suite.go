package types

type Suite struct {
	name   string
	weight int
	tasks  []Task
}

func (s Suite) Name() string {
	return s.name
}

func (s Suite) Weight() int {
	return s.weight
}

// Tasks returns the task by scope. If scope is empty, it returns all
func (s Suite) Tasks(scope string) []Task {
	if scope == "" {
		return s.tasks
	}

	var tasks []Task
	for _, task := range s.tasks {
		if task.Scope() == scope {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func NewSuite(name string, weight int, tasks []Task) Suite {
	return Suite{name: name, weight: weight, tasks: tasks}
}
