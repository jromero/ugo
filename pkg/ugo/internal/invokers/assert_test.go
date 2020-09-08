package invokers_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"

	"github.com/jromero/ugo/pkg/ugo/internal/invokers"
	"github.com/jromero/ugo/pkg/ugo/internal/tasks"
	"github.com/jromero/ugo/pkg/ugo/types"
)

func TestAssertInvoker(t *testing.T) {
	spec.Run(t, "AssertInvoker", func(t *testing.T, when spec.G, it spec.S) {
		when("ignoreLines is '...'", func() {
			it("ignores specific lines", func() {
				invoker := invokers.AssertInvoker{}
				assertContainsTask := tasks.NewAssertContainsTask(types.ScopeDefault, `line 1
...
...line 4
line 5...
...
`,
					tasks.WithIgnoreLines("..."))

				assert.True(t, invoker.Supports(assertContainsTask))
				_, err := invoker.Invoke(assertContainsTask, "", `line 1
line 2
line 3
...line 4
line 5...
line 6
line 7`)
				assert.Nil(t, err)
			})

			it("doesn't ignore lines that contain those chars", func() {
				invoker := invokers.AssertInvoker{}
				assertContainsTask := tasks.NewAssertContainsTask(types.ScopeDefault, `line 1
...line 2
line 3...
`,
					tasks.WithIgnoreLines("..."))

				assert.True(t, invoker.Supports(assertContainsTask))
				_, err := invoker.Invoke(assertContainsTask, "", `line 1
line 2
line 3`)

				assert.EqualError(t, err, "output didn't match:\n(?ms).*line 1\n...line 2\nline 3...\n.*")
			})
		})
	})
}
