package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/thatisuday/commando"

	"github.com/jromero/ugo/pkg/ugo"
	"github.com/jromero/ugo/pkg/ugo/types"
)

type cliArgs = map[string]commando.ArgValue
type cliFlags = map[string]commando.FlagValue

func main() {
	commando.
		SetExecutableName("ugo").
		SetVersion(ugo.Version).
		SetDescription("Ugo helps you execute and test your tutorials.")

	commando.
		Register("run").
		SetShortDescription("run through some tutorials").
		SetDescription(`Run through some tutorials.`).
		AddFlag("path,p", "path to your tutorials", commando.String, ".").
		AddFlag("recursive,r,R", "recursively look for tutorials", commando.Bool, false).
		AddFlag("verbose,v", "verbose output", commando.Bool, false).
		AddFlag("extensions,e", "comma separated list of tutorial file extensions", commando.String, ".md,.mdx").
		AddFlag("keep,k", "keep the generated files", commando.Bool, false).
		SetAction(func(args cliArgs, flags cliFlags) {
			logger := &ugo.Logger{
				Level: ugo.INFO,
			}

			path, err := flags["path"].GetString()
			if err != nil {
				fatalError(logger, err, 1)
			}

			recursive, err := flags["recursive"].GetBool()
			if err != nil {
				fatalError(logger, err, 1)
			}

			verbose, err := flags["verbose"].GetBool()
			if err != nil {
				fatalError(logger, err, 1)
			}

			if verbose {
				logger.Level = ugo.DEBUG
			}

			extList, err := flags["extensions"].GetString()
			if err != nil {
				fatalError(logger, err, 1)
			}
			extensions := strings.Split(extList, ",")

			keep, err := flags["keep"].GetBool()
			if err != nil {
				fatalError(logger, err, 1)
			}

			files, err := searchForFiles(path, recursive, extensions)
			if err != nil {
				fatalError(logger, err, 1)
			}

			var plans []types.Plan
			for _, file := range files {
				logger.Debug("Reading file: %s", file)

				content, err := ioutil.ReadFile(file)
				if err != nil {
					fatalError(logger, err, 2)
				}

				p, err := ugo.Parse(string(content))
				if err != nil {
					if err == ugo.NoSuiteError {
						continue
					}

					fatalError(logger, fmt.Errorf("parsing '%s': %s", file, err.Error()), 2)
				}

				plans = append(plans, p)
			}

			if len(plans) == 0 {
				logger.Info("Nothing found to execute or test.")
				return
			}

			_, err = ugo.Invoke(logger, keep, ugo.Aggregate(plans...))
			if err != nil {
				fatalError(logger, err, 3)
			}
			logger.Info("Nothing broken. Good job!")
		})

	commando.Parse(nil)
}

func searchForFiles(path string, recursive bool, extensions []string) (files []string, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		if f, err := ioutil.ReadDir(path); err != nil {
			return nil, err
		} else {
			for _, fileInfo := range f {
				ext := filepath.Ext(fileInfo.Name())
				hidden := strings.HasPrefix(filepath.Base(fileInfo.Name()), ".")
				if !fileInfo.IsDir() && lo.Contains(extensions, ext) && !hidden {
					files = append(files, filepath.Join(path, fileInfo.Name()))
				} else if recursive {
					f, err := searchForFiles(filepath.Join(path, fileInfo.Name()), recursive, extensions)
					if err != nil {
						return nil, err
					}
					files = append(files, f...)
				}
			}
		}
	} else {
		ext := filepath.Ext(fileInfo.Name())
		hidden := strings.HasPrefix(filepath.Base(fileInfo.Name()), ".")
		if lo.Contains(extensions, ext) && !hidden {
			files = append(files, path)
		}
	}

	for i := 0; i < len(files); i++ {
		files[i], err = filepath.Abs(files[i])
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func fatalError(logger types.Logger, err error, code int) {
	logger.Error(err)
	os.Exit(code)
}
