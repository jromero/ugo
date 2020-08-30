package ugo

import "sort"

func Aggregate(plans ...Plan) Plan {
	var suites []Suite
	for _, plan := range plans {
		suites = append(suites, plan.suites...)
	}

	return *NewPlan(aggregateSuites(suites))
}

func aggregateSuites(suites []Suite) []Suite {
	var newSuites []Suite
	groupedSuites := map[string][]Suite{}
	for _, s := range suites {
		group, ok := groupedSuites[s.name]
		if !ok {
			group = []Suite{}
		}

		groupedSuites[s.name] = append(group, s)
	}

	for _, ss := range groupedSuites {
		if len(ss) == 1 {
			newSuites = append(newSuites, ss[0])
			continue
		}

		// keep first instance values
		name := ss[0].name
		weight := ss[0].weight

		// aggregate tasks by weight
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].weight < ss[j].weight
		})

		tasks := ss[0].tasks
		for i := 1; i < len(ss); i++ {
			tasks = append(tasks, ss[i].tasks...)
		}

		newSuites = append(newSuites, *NewSuite(name, weight, tasks))
	}

	sort.Slice(newSuites, func(i, j int) bool {
		return newSuites[i].weight < newSuites[j].weight
	})

	return newSuites
}
