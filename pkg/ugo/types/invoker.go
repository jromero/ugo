package types

type Invoker interface {
	Supports(task Task) bool
	Invoke(task Task, keep bool, workDir, prevOutput string) (output string, err error)
}
