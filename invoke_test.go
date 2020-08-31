package ugo_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"

	"github.com/jromero/ugo"
	"github.com/jromero/ugo/internal/tasks"
	"github.com/jromero/ugo/internal/types"
)

func TestInvoke(t *testing.T) {
	spec.Run(t, "#Invoke", func(t *testing.T, when spec.G, it spec.S) {
		when("exec", func() {
			when("exit code doesn't match expected", func() {
				it("errors", func() {
					err := ugo.Invoke(*ugo.NewPlan([]ugo.Suite{
						*ugo.NewSuite("test1", 0, []types.Task{
							tasks.NewExecTask("exit 1", 0),
						}),
					}))

					assert.Error(t, err, "exec task #1 failed: exit status 1")
				})
			})

			when("exit code matches expected", func() {
				it("doesn't error", func() {
					err := ugo.Invoke(*ugo.NewPlan([]ugo.Suite{
						*ugo.NewSuite("test1", 0, []types.Task{
							tasks.NewExecTask("exit 1", 1),
						}),
					}))

					assert.Nil(t, err)
				})
			})

			when("exit code expected is set to -1", func() {
				it("doesn't error", func() {
					err := ugo.Invoke(*ugo.NewPlan([]ugo.Suite{
						*ugo.NewSuite("test1", 0, []types.Task{
							tasks.NewExecTask("exit 99", -1),
						}),
					}))

					assert.Nil(t, err)
				})
			})
		})

		when("assert", func() {
			when("contains", func() {
				when("matches content", func() {
					it("doesn't error", func() {
						err := ugo.Invoke(*ugo.NewPlan([]ugo.Suite{
							*ugo.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask(`
echo "hello #1"
echo "hello #2"
`, 0),
								tasks.NewAssertContainsTask("hello #1\nhello #2"),
							}),
						}))

						assert.Nil(t, err)
					})
				})

				when("multiple consecutive asserts", func() {
					it("searches in output", func() {
						err := ugo.Invoke(*ugo.NewPlan([]ugo.Suite{
							*ugo.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask("echo hello1;echo hello2", 0),
								tasks.NewAssertContainsTask("hello1"),
								tasks.NewAssertContainsTask("hello2"),
							}),
						}))

						assert.Nil(t, err)
					})
				})

				when("content has ansi codes", func() {
					it("matches ignoring ansi", func() {
						err := ugo.Invoke(*ugo.NewPlan([]ugo.Suite{
							*ugo.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask(`
echo -e "\x1b[38;5;140mfoo\x1b[0mbar"
`, 0),
								tasks.NewAssertContainsTask("foobar"),
							}),
						}))

						assert.Nil(t, err)
					})
				})

				when("contents don't contain", func() {
					it("errors", func() {
						err := ugo.Invoke(*ugo.NewPlan([]ugo.Suite{
							*ugo.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask(`echo "hello #1"`, 0),
								tasks.NewAssertContainsTask("hello #2"),
							}),
						}))

						assert.EqualError(t, err, "no output contained:\nhello #2")
					})
				})
			})
		})
	})
}
