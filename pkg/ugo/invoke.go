package ugo

import (
	"fmt"
	"io/ioutil"

	"github.com/jromero/ugo/pkg/ugo/internal/invokers"
	"github.com/jromero/ugo/pkg/ugo/types"
)

func Invoke(logger types.Logger, plan types.Plan) error {
	taskInvokers := []types.Invoker{
		invokers.NewAssertInvoker(logger),
		invokers.NewExecInvoker(logger),
		invokers.NewFileInvoker(logger),
	}

	for _, suite := range plan.Suites() {
		logger.Info("Suite '%s' executing...", suite.Name())
		logger.AddBreadcrumb(suite.Name())

		workDir, err := ioutil.TempDir("", fmt.Sprintf("suite-%s-*", suite.Name()))
		if err != nil {
			return err
		}
		logger.Debug("Working directory: %s", workDir)

		orderedTasks := append(append(
			suite.Tasks(types.ScopeSetup),
			suite.Tasks(types.ScopeDefault)...),
			suite.Tasks(types.ScopeTeardown)...)

		var aggrOutput string
		for i, task := range orderedTasks {
			logger.Info("Running task #%d", i+1)
			taskName := fmt.Sprintf("#%d-%s:%s", i+1, task.Scope(), task.Name())
			logger.AddBreadcrumb(taskName)

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

			logger.PopBreadcrumb(taskName)
		}

		logger.PopBreadcrumb(suite.Name())
	}

	return nil
}
