package pkg_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"

	"github.com/jromero/ugo/pkg"
	"github.com/jromero/ugo/pkg/internal/tasks"
	"github.com/jromero/ugo/pkg/types"
)

func TestAggregate(t *testing.T) {
	spec.Run(t, "#Aggregate", func(t *testing.T, when spec.G, it spec.S) {
		when("cross-file-tutorial", func() {
			it.Focus("should rely on weight, not filename", func() {
				tutorialDir := filepath.Join("testdata", "cross-file-tutorial")

				fileInfos, err := ioutil.ReadDir(tutorialDir)
				assert.Nil(t, err)

				var plans []types.Plan
				for _, fileInfo := range fileInfos {
					if !fileInfo.IsDir() {
						bytes, err := ioutil.ReadFile(filepath.Join(tutorialDir, fileInfo.Name()))
						assert.Nil(t, err)

						p, err := pkg.Parse(string(bytes))
						assert.Nil(t, err)

						plans = append(plans, p)
					}
				}

				plan := pkg.Aggregate(plans...)
				assert.Equal(t, types.NewPlan([]types.Suite{types.NewSuite("cross-file-tutorial", 1, []types.Task{
					tasks.NewExecTask("echo hello 1", 0),
					tasks.NewExecTask("echo hello 2", 0),
					tasks.NewExecTask("echo hello 3", 0),
				})}), plan)
			})
		})
	})
}
