package types

type Plan struct {
	suites []Suite
}

func NewPlan(suites []Suite) Plan {
	return Plan{suites: suites}
}

func (p Plan) Suites() []Suite {
	return p.suites
}
