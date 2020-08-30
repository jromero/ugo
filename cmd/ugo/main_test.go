package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	spec.Run(t, "#searchForFiles", func(t *testing.T, when spec.G, it spec.S) {
		when("path is directory", func() {
			when("recursive=true", func() {
				it("searches recursively", func() {
					cwd, err := os.Getwd()
					assert.Nil(t, err)

					files, err := searchForFiles(filepath.Join("testdata", "tutorials"), true)
					assert.Nil(t, err)
					assert.Equal(t, []string{
						filepath.Join(cwd, "testdata", "tutorials", "nested-1", "nested-1a", "tutorial-nested-1a.md"),
						filepath.Join(cwd, "testdata", "tutorials", "nested-1", "nested-1b", "tutorial-nested-1b-1.md"),
						filepath.Join(cwd, "testdata", "tutorials", "nested-1", "nested-1b", "tutorial-nested-1b-2.md"),
						filepath.Join(cwd, "testdata", "tutorials", "nested-1", "tutorial-nested-1.md"),
						filepath.Join(cwd, "testdata", "tutorials", "nested-2", "tutorial-nested-2.md"),
						filepath.Join(cwd, "testdata", "tutorials", "tutorial-1.md"),
						filepath.Join(cwd, "testdata", "tutorials", "tutorial-2.md"),
					}, files)
				})
			})

			when("--recursive=false", func() {
				it("doesn't traverse directories", func() {
					cwd, err := os.Getwd()
					assert.Nil(t, err)

					files, err := searchForFiles(filepath.Join("testdata", "tutorials"), false)
					assert.Nil(t, err)
					assert.Equal(t, []string{
						filepath.Join(cwd, "testdata", "tutorials", "tutorial-1.md"),
						filepath.Join(cwd, "testdata", "tutorials", "tutorial-2.md"),
					}, files)
				})
			})
		})

		when("path is a file", func() {
			it("returns file if exists", func() {
				cwd, err := os.Getwd()
				assert.Nil(t, err)

				files, err := searchForFiles(filepath.Join("testdata", "tutorials", "tutorial-1.md"), false)
				assert.Nil(t, err)
				assert.Equal(t, []string{
					filepath.Join(cwd, "testdata", "tutorials", "tutorial-1.md"),
				}, files)
			})

			it("errors if file doesn't exist", func() {
				_, err := searchForFiles(filepath.Join("testdata", "tutorials", "non-exists.md"), false)
				assert.EqualError(t, err, "stat testdata/tutorials/non-exists.md: no such file or directory")
			})
		})
	})
}
