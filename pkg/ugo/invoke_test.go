package ugo_test

import (
	"os"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"

	"github.com/jromero/ugo/pkg/ugo"
	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

func TestInvoke(t *testing.T) {
	logger := &ugo.Logger{}

	spec.Run(t, "#Invoke", func(t *testing.T, when spec.G, it spec.S) {
		when("exec", func() {
			when("exit code doesn't match expected", func() {
				it("errors", func() {
					wds, err := ugo.Invoke(
						logger,
						false,
						types.NewPlan([]types.Suite{
							types.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask(types.ScopeDefault, "exit 1", 0),
							}),
						}))

					assert.Error(t, err, "exec task #1 failed: exit status 1")
					assert.Empty(t, wds)
				})
			})

			when("exit code matches expected", func() {
				it("doesn't error", func() {
					wds, err := ugo.Invoke(
						logger,
						false,
						types.NewPlan([]types.Suite{
							types.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask(types.ScopeDefault, "exit 1", 1),
							}),
						}))

					assert.Nil(t, err)
					assert.Empty(t, wds)
				})
			})

			when("exit code expected is set to -1", func() {
				it("doesn't error", func() {
					wds, err := ugo.Invoke(
						logger,
						false,
						types.NewPlan([]types.Suite{
							types.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask(types.ScopeDefault, "exit 99", -1),
							}),
						}))

					assert.Nil(t, err)
					assert.Empty(t, wds)
				})
			})

			when("asked to keep generated files", func() {
				it("keeps scripts", func() {
					wds, err := ugo.Invoke(
						logger,
						true,
						types.NewPlan([]types.Suite{
							types.NewSuite("test1", 0, []types.Task{
								tasks.NewExecTask(types.ScopeDefault, "exit 0", 0),
							}),
						}))

					assert.Nil(t, err)
					assert.Len(t, wds, 1)
					entry, err := os.ReadDir(wds[0])
					assert.Nil(t, err)
					assert.Len(t, entry, 1)
					os.RemoveAll(wds[0])
				})
			})
		})

		when("assert", func() {
			when("contains", func() {
				when("matches content", func() {
					it("doesn't error", func() {
						wds, err := ugo.Invoke(
							logger,
							false,
							types.NewPlan([]types.Suite{
								types.NewSuite("test1", 0, []types.Task{
									tasks.NewExecTask(types.ScopeDefault, `
echo "hello #1"
echo "hello #2"
`, 0),
									tasks.NewAssertContainsTask(types.ScopeDefault, "hello #1\nhello #2"),
								}),
							}))

						assert.Nil(t, err)
						assert.Empty(t, wds)
					})
				})

				when("multiple consecutive asserts", func() {
					it("searches in output", func() {
						wds, err := ugo.Invoke(
							logger,
							false,
							types.NewPlan([]types.Suite{
								types.NewSuite("test1", 0, []types.Task{
									tasks.NewExecTask(types.ScopeDefault, "echo hello1;echo hello2", 0),
									tasks.NewAssertContainsTask(types.ScopeDefault, "hello1"),
									tasks.NewAssertContainsTask(types.ScopeDefault, "hello2"),
								}),
							}))

						assert.Nil(t, err)
						assert.Empty(t, wds)
					})
				})

				when("content has ansi codes", func() {
					it("matches ignoring ansi", func() {
						wds, err := ugo.Invoke(
							logger,
							false,
							types.NewPlan([]types.Suite{
								types.NewSuite("test1", 0, []types.Task{
									tasks.NewExecTask(types.ScopeDefault, `
echo -e "\x1b[38;5;140mfoo\x1b[0mbar"
`, 0),
									tasks.NewAssertContainsTask(types.ScopeDefault, "foobar"),
								}),
							}))

						assert.Nil(t, err)
						assert.Empty(t, wds)
					})
				})

				when("contents don't contain", func() {
					it("errors", func() {
						wds, err := ugo.Invoke(
							logger,
							false,
							types.NewPlan([]types.Suite{
								types.NewSuite("test1", 0, []types.Task{
									tasks.NewExecTask(types.ScopeDefault, `echo "hello #1"`, 0),
									tasks.NewAssertContainsTask(types.ScopeDefault, "hello #2"),
								}),
							}))

						assert.EqualError(t, err, "Output did not contain:\nhello #2")
						assert.Empty(t, wds)
					})
				})
			})
		})

		when("scopes", func() {
			it("executes in order: setup, default, teardown", func() {
				wds, err := ugo.Invoke(
					logger,
					false,
					types.NewPlan([]types.Suite{
						types.NewSuite("test1", 0, []types.Task{
							tasks.NewExecTask(types.ScopeSetup, `echo hello > test-file.txt`, 0),
							tasks.NewExecTask(types.ScopeTeardown, `rm test-file.txt`, 0),
							tasks.NewExecTask(types.ScopeTeardown, `cat test-file.txt`, 1),
							tasks.NewExecTask(types.ScopeDefault, `cat test-file.txt`, 0),
						}),
					}),
				)

				assert.Nil(t, err)
				assert.Empty(t, wds)
			})

		})
	})
}
