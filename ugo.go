package ugo

type Plan struct {
	suites []Suite
}

func NewPlan(suites []Suite) *Plan {
	return &Plan{suites: suites}
}

type Suite struct {
	name   string
	weight int
	tasks  []Task
}

func NewSuite(name string, weight int, tasks []Task) *Suite {
	return &Suite{name: name, weight: weight, tasks: tasks}
}

const (
	TypeFile           = "file"
	TypeExec           = "exec"
	TypeAssertContains = "assert"
	AttrFilename       = "filename"
	AttrExitCode       = "exit-code"
	AttrAssertion      = "assertion"
)

type Task struct {
	typ      string
	contents string
	attr     map[string]interface{}
}

func (t *Task) Type() string {
	return t.typ
}

// NewExecTask creates a Task of TypeExec with contents of script and expected exit code.
func NewExecTask(contents string, exitCode int) *Task {
	return &Task{
		typ:      TypeExec,
		contents: contents,
		attr: map[string]interface{}{
			AttrExitCode: exitCode,
		},
	}
}

// NewFileTask creates a Task of TypeFile with filename and contents of file.
func NewFileTask(filename, contents string) *Task {
	return &Task{
		typ:      TypeFile,
		contents: contents,
		attr: map[string]interface{}{
			AttrFilename: filename,
		},
	}
}

// NewAssertContainsTask creates a Task of TypeAssertContains.
func NewAssertContainsTask(contents string) *Task {
	return &Task{
		typ:      TypeAssertContains,
		contents: contents,
	}
}
