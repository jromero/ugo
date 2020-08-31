package types

type Invoker interface {
	Supports(task Task) bool
	Invoke(task Task, workDir, prevOutput string) (output string, err error)
}
