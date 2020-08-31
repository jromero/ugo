package ugo

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jromero/ugo/pkg/ugo/internal/invokers"
	"github.com/jromero/ugo/pkg/ugo/types"
)

var taskInvokers = []types.Invoker{
	&invokers.AssertInvoker{},
	&invokers.ExecInvoker{},
	&invokers.FileInvoker{},
}

func Invoke(plan types.Plan) error {
	prevFlags := log.Flags()
	defer func() { log.SetFlags(prevFlags) }()
	prevPrefix := log.Prefix()
	defer func() { log.SetPrefix(prevPrefix) }()

	log.SetFlags(0)

	for _, suite := range plan.Suites() {
		log.SetPrefix(fmt.Sprintf("[%s] ", suite.Name()))
		log.Printf("Suite '%s' executing...", suite.Name())

		workDir, err := ioutil.TempDir("", fmt.Sprintf("suite-%s-*", suite.Name()))
		if err != nil {
			return err
		}
		log.Println("Working directory:", workDir)

		orderedTasks := append(append(
			suite.Tasks(types.ScopeSetup),
			suite.Tasks(types.ScopeDefault)...),
			suite.Tasks(types.ScopeTeardown)...)

		var aggrOutput string
		for i, task := range orderedTasks {
			log.SetPrefix(fmt.Sprintf("[%s][#%d-%s:%s] ", suite.Name(), i+1, task.Scope(), task.Name()))
			log.Printf("--> Running task #%d", i+1)

			for _, invoker := range taskInvokers {
				if invoker.Supports(task) {
					output, err := invoker.Invoke(task, workDir, aggrOutput)
					if err != nil {
						return err
					}

					aggrOutput += output
					break
				}
			}
		}
	}

	return nil
}
