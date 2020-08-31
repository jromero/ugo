package types

const (
	ScopeDefault  = "default"
	ScopeSetup    = "setup"
	ScopeTeardown = "teardown"
)

type Task interface {
	Name() string

	// the scope of this task, one of ScopeSetup, ScopeDefault, ScopeTeardown
	Scope() string
}
