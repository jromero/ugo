package ugo_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"

	"github.com/jromero/ugo/pkg/ugo"
	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

func TestParse(t *testing.T) {
	spec.Run(t, "#Parse", func(t *testing.T, when spec.G, it spec.S) {

		when("suite", func() {
			it("parses suite", func() {
				plan, err := ugo.Parse(`
<!-- test:suite=test1;weight=1 -->
`)
				assert.Nil(t, err)
				assert.Equal(t, types.NewPlan([]types.Suite{types.NewSuite("test1", 1, nil)}), plan)
			})

			it("no suite found", func() {
				_, err := ugo.Parse(`some content without a suite definition`)
				assert.EqualError(t, err, "no suite found")
			})

			it("multiple suites (are ordered)", func() {
				plan, err :=
					ugo.Parse(`
<!-- test:suite=test1;weight=1 -->
<!-- test:exec -->
` + "```shell bash" + `
echo "hello #1"!
` + "```" + `
<!-- test:suite=test2;weight=3 -->
<!-- test:exec -->
` + "```shell bash" + `
echo "hello #2"!
` + "```" + `
<!-- test:suite=test3;weight=2 -->
<!-- test:exec -->
` + "```shell bash" + `
echo "hello #3"!
` + "```" + `
`)
				assert.Nil(t, err)

				expected := types.NewPlan([]types.Suite{
					types.NewSuite("test1", 1, []types.Task{
						tasks.NewExecTask(`echo "hello #1"!`, 0),
					}),
					types.NewSuite("test3", 2, []types.Task{
						tasks.NewExecTask(`echo "hello #3"!`, 0),
					}),
					types.NewSuite("test2", 3, []types.Task{
						tasks.NewExecTask(`echo "hello #2"!`, 0),
					}),
				})

				assert.Equal(t, expected, plan)
			})

			it("multiple suites (aggregated based on weight)", func() {
				plan, err :=
					ugo.Parse(`
<!-- test:suite=test1;weight=2 -->
<!-- test:exec -->
` + "```shell bash" + `
echo "hello #1"!
` + "```" + `
<!-- test:suite=test2;weight=3 -->
<!-- test:exec -->
` + "```shell bash" + `
echo "hello #2"!
` + "```" + `
<!-- test:suite=test1;weight=1 -->
<!-- test:exec -->
` + "```shell bash" + `
echo "hello #3"!
` + "```" + `
`)
				assert.Nil(t, err)

				expected := types.NewPlan([]types.Suite{
					types.NewSuite("test1", 2, []types.Task{
						tasks.NewExecTask(`echo "hello #3"!`, 0),
						tasks.NewExecTask(`echo "hello #1"!`, 0),
					}),
					types.NewSuite("test2", 3, []types.Task{
						tasks.NewExecTask(`echo "hello #2"!`, 0),
					}),
				})
				assert.Equal(t, expected, plan)
			})
		})

		when("tasks", func() {
			when("unknown task type", func() {
				it("errors", func() {
					_, err := ugo.Parse(`
<!-- test:suite=test1 -->

<!-- test:some-other-value -->
` + "```shell bash" + `
echo "hello #1"!
` + "```" + `
`)
					assert.EqualError(t, err, "unknown task 'some-other-value'")
				})
			})

			when("file", func() {
				it("file task is parsed", func() {
					plan, err := ugo.Parse(`
<!-- test:suite=test1;weight=1 -->
<!-- test:file=some-file -->
` + "```" + `
some-content
` + "```" + `
`)
					assert.Nil(t, err)

					expected := types.NewPlan([]types.Suite{types.NewSuite("test1", 1, []types.Task{
						tasks.NewFileTask("some-file", "some-content"),
					})})
					assert.Equal(t, expected, plan)
				})
			})

			when("exec", func() {
				it("is parsed", func() {
					plan, err := ugo.Parse(`
<!-- test:suite=test1 -->
	
<!-- test:exec -->
` + "```shell bash" + `
echo "hello"!
` + "```" + `
`)
					assert.Nil(t, err)

					expected := types.NewPlan([]types.Suite{types.NewSuite("test1", 0, []types.Task{
						tasks.NewExecTask(`echo "hello"!`, 0),
					})})
					assert.Equal(t, expected, plan)
				})

				it("parses exit code", func() {
					plan, err := ugo.Parse(`
<!-- test:suite=test1 -->
	
<!-- test:exec;exit-code=1 -->
` + "```shell bash" + `
exit 1
` + "```" + `
`)
					assert.Nil(t, err)

					expected := types.NewPlan([]types.Suite{types.NewSuite("test1", 0, []types.Task{
						tasks.NewExecTask(`exit 1`, 1),
					})})
					assert.Equal(t, expected, plan)
				})

				it("parses exit code (-1)", func() {
					plan, err := ugo.Parse(`
<!-- test:suite=test1 -->
	
<!-- test:exec;exit-code=-1 -->
` + "```shell bash" + `
exit 99
` + "```" + `
`)
					assert.Nil(t, err)

					expected := types.NewPlan([]types.Suite{types.NewSuite("test1", 0, []types.Task{
						tasks.NewExecTask(`exit 99`, -1),
					})})
					assert.Equal(t, expected, plan)
				})
			})

			when("assertion", func() {
				when("contains", func() {
					it("is parsed", func() {
						plan, err := ugo.Parse(`
<!-- test:suite=test1 -->
	
<!-- test:assert=contains -->
` + "```text" + `
some-output
` + "```" + `
`)
						assert.Nil(t, err)

						expected := types.NewPlan([]types.Suite{types.NewSuite("test1", 0, []types.Task{
							tasks.NewAssertContainsTask(`some-output`),
						})})
						assert.Equal(t, expected, plan)
					})
				})
			})

			it("multiple tasks", func() {
				plan, err :=
					ugo.Parse(`
<!-- test:suite=test1 -->

<!-- test:file=some-file -->
` + "```toml" + `
[root]
key=value
` + "```" + `
	
<!-- test:exec -->
` + "```shell bash" + `
echo "hello"!
` + "```" + `
`)
				assert.Nil(t, err)

				expected := types.NewPlan([]types.Suite{types.NewSuite("test1", 0, []types.Task{
					tasks.NewFileTask("some-file", "[root]\nkey=value"),
					tasks.NewExecTask(`echo "hello"!`, 0),
				})})

				assert.Equal(t, expected, plan)
			})

			it("multiple code blocks after task", func() {
				plan, err := ugo.Parse(`
<!-- test:suite=test1 -->
	
<!-- test:exec -->
` + "```shell bash" + `
echo "hello #1"!
` + "```" + `

` + "```shell bash" + `
echo "hello #2"!
` + "```" + `
`)
				assert.Nil(t, err)

				expected := types.NewPlan([]types.Suite{types.NewSuite("test1", 0, []types.Task{
					tasks.NewExecTask(`echo "hello #1"!`, 0),
				})})
				assert.Equal(t, expected, plan)
			})
		})
	})
}
