package ugo_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"

	"github.com/jromero/ugo"
)

func TestExecute(t *testing.T) {
	spec.Run(t, "#Execute", func(t *testing.T, when spec.G, it spec.S) {
		when("exec", func() {
			when("exit code doesn't match expected", func() {
				it("errors", func() {
					err := ugo.Execute(*ugo.NewPlan([]ugo.Suite{
						*ugo.NewSuite("test1", 0, []ugo.Task{
							*ugo.NewExecTask("exit 1", 0),
						}),
					}))

					assert.Error(t, err, "exec task #1 failed: exit status 1")
				})
			})

			when("exit code matches expected", func() {
				it("doesn't error", func() {
					err := ugo.Execute(*ugo.NewPlan([]ugo.Suite{
						*ugo.NewSuite("test1", 0, []ugo.Task{
							*ugo.NewExecTask("exit 1", 1),
						}),
					}))

					assert.Nil(t, err)
				})
			})

			when("exit code expected is set to -1", func() {
				it("doesn't error", func() {
					err := ugo.Execute(*ugo.NewPlan([]ugo.Suite{
						*ugo.NewSuite("test1", 0, []ugo.Task{
							*ugo.NewExecTask("exit 99", -1),
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
						err := ugo.Execute(*ugo.NewPlan([]ugo.Suite{
							*ugo.NewSuite("test1", 0, []ugo.Task{
								*ugo.NewExecTask(`
echo "hello #1"
echo "hello #2"
`, 0),
								*ugo.NewAssertContainsTask("hello #1\nhello #2"),
							}),
						}))

						assert.Nil(t, err)
					})
				})
				
				when("content has ansi codes", func() {
					it("matches ignoring ansi", func() {
						err := ugo.Execute(*ugo.NewPlan([]ugo.Suite{
							*ugo.NewSuite("test1", 0, []ugo.Task{
								*ugo.NewExecTask(`
echo -e "\x1b[38;5;140mfoo\x1b[0mbar"
`, 0),
								*ugo.NewAssertContainsTask("foobar"),
							}),
						}))

						assert.Nil(t, err)
					})
				})

				when("contents don't contain", func() {
					it("doesn't error", func() {
						err := ugo.Execute(*ugo.NewPlan([]ugo.Suite{
							*ugo.NewSuite("test1", 0, []ugo.Task{
								*ugo.NewExecTask(`echo "hello #1"`, 0),
								*ugo.NewAssertContainsTask("hello #2"),
							}),
						}))

						assert.EqualError(t, err, "assert task #2 failed: no output contained:\nhello #2")
					})
				})
			})
		})
	})
}
