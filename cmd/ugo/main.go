package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/thatisuday/commando"

	"github.com/jromero/ugo"
)

type cliArgs = map[string]commando.ArgValue
type cliFlags = map[string]commando.FlagValue

func main() {
	commando.
		SetExecutableName("ugo").
		SetVersion("0.0.1").
		SetDescription("Ugo helps you execute and test your tutorials.")

	commando.
		Register("run").
		SetShortDescription("run through some tutorials").
		SetDescription(`Run through some tutorials.`).
		AddFlag("path,p", "path to your tutorials", commando.String, ".").
		AddFlag("recursive,r,R", "recursively look for tutorials", commando.Bool, false).
		AddFlag("verbose,v", "verbose output", commando.Bool, false).
		SetAction(func(args cliArgs, flags cliFlags) {
			path, err := flags["path"].GetString()
			if err != nil {
				fatalError(1, err.Error())
			}

			recursive, err := flags["recursive"].GetBool()
			if err != nil {
				fatalError(1, err.Error())
			}

			verbose, err := flags["verbose"].GetBool()
			if err != nil {
				fatalError(1, err.Error())
			}

			files, err := searchForFiles(path, recursive)
			if err != nil {
				fatalError(1, err.Error())
			}

			var plans []ugo.Plan
			for _, file := range files {
				if verbose {
					println("reading file:", file)
				}

				content, err := ioutil.ReadFile(file)
				if err != nil {
					fatalError(2, err.Error())
				}

				p, err := ugo.Parse(string(content))
				if err != nil {
					if err == ugo.NoSuiteError {
						continue
					}

					fatalError(2, "parsing '%s': %s", file, err.Error())
				}

				plans = append(plans, p)
			}

			if len(plans) == 0 {
				println("nothing found to execute or test")
				return
			}

			err = ugo.Execute(ugo.Aggregate(plans...))
			if err != nil {
				fatalError(3, err.Error())
			}

			println("Nothing broken. Good job!")
		})

	commando.Parse(nil)
}

func searchForFiles(path string, recursive bool) (files []string, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		if f, err := ioutil.ReadDir(path); err != nil {
			return nil, err
		} else {
			for _, fileInfo := range f {
				if !fileInfo.IsDir() {
					files = append(files, filepath.Join(path, fileInfo.Name()))
				} else if recursive {
					f, err := searchForFiles(filepath.Join(path, fileInfo.Name()), recursive)
					if err != nil {
						return nil, err
					}
					files = append(files, f...)
				}
			}
		}
	} else {
		files = append(files, path)
	}

	for i := 0; i < len(files); i++ {
		files[i], err = filepath.Abs(files[i])
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func fatalError(code int, msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf("Error: "+msg, args...)
	} else {
		fmt.Println("Error: " + msg)
	}
	os.Exit(code)
}
