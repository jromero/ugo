package types

type Parser interface {
	AttemptParse(scope, taskDefinition, nextCodeBlock string) (Task, error)
}
