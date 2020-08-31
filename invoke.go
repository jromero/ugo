package ugo

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jromero/ugo/internal/invokers"
	"github.com/jromero/ugo/internal/types"
)

var taskInvokers = []types.Invoker{
	&invokers.AssertInvoker{},
	&invokers.ExecInvoker{},
	&invokers.FileInvoker{},
}

func Invoke(plan Plan) error {
	prevFlags := log.Flags()
	defer func() { log.SetFlags(prevFlags) }()
	prevPrefix := log.Prefix()
	defer func() { log.SetPrefix(prevPrefix) }()

	log.SetFlags(0)

	for _, suite := range plan.suites {
		log.SetPrefix(fmt.Sprintf("[%s] ", suite.name))
		log.Printf("Suite '%s' executing...", suite.name)

		workDir, err := ioutil.TempDir("", fmt.Sprintf("suite-%s-*", suite.name))
		if err != nil {
			return err
		}
		log.Println("Working directory:", workDir)

		var aggrOutput string
		for i, task := range suite.tasks {
			log.SetPrefix(fmt.Sprintf("[%s][task#%d] ", suite.name, i+1))
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
