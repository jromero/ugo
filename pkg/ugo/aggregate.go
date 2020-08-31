package ugo

import (
	"sort"

	"github.com/jromero/ugo/pkg/ugo/types"
)

func Aggregate(plans ...types.Plan) types.Plan {
	var suites []types.Suite
	for _, plan := range plans {
		suites = append(suites, plan.Suites()...)
	}

	return types.NewPlan(aggregateSuites(suites))
}

func aggregateSuites(suites []types.Suite) []types.Suite {
	var newSuites []types.Suite
	groupedSuites := map[string][]types.Suite{}
	for _, s := range suites {
		group, ok := groupedSuites[s.Name()]
		if !ok {
			group = []types.Suite{}
		}

		groupedSuites[s.Name()] = append(group, s)
	}

	for _, ss := range groupedSuites {
		if len(ss) == 1 {
			newSuites = append(newSuites, ss[0])
			continue
		}

		// keep first instance values
		name := ss[0].Name()
		weight := ss[0].Weight()

		// aggregate tasks by weight
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Weight() < ss[j].Weight()
		})

		tasks := ss[0].Tasks()
		for i := 1; i < len(ss); i++ {
			tasks = append(tasks, ss[i].Tasks()...)
		}

		newSuites = append(newSuites, types.NewSuite(name, weight, tasks))
	}

	sort.Slice(newSuites, func(i, j int) bool {
		return newSuites[i].Weight() < newSuites[j].Weight()
	})

	return newSuites
}
