/*
Copyright 2023 Go Imports Organizer Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"

	"github.com/go-imports-organizer/gio/pkg/config"
	"github.com/go-imports-organizer/gio/pkg/excludes"
	"github.com/go-imports-organizer/gio/pkg/groups"
	"github.com/go-imports-organizer/gio/pkg/imports"
	"github.com/go-imports-organizer/gio/pkg/module"
	"github.com/go-imports-organizer/gio/pkg/version"
)

var (
	wg    sync.WaitGroup
	files = make(chan string)
)

func main() {
	// If the -l flag is set, only return a list of what files would have been formatted, but don't make any changes
	listOnly := flag.Bool("l", false, "list files that need to be organized (no changes made)")
	versionOnly := flag.Bool("v", false, "print version and exit")
	flag.Parse()

	// set CPUPROFILE=<filename> to create a <filename>.pprof cpu profile file
	if len(os.Getenv("CPUPROFILE")) != 0 {
		f, err := os.Create(os.Getenv("CPUPROFILE"))
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *versionOnly {
		fmt.Fprintf(os.Stdout, "version %s\n", version.Get())
		os.Exit(0)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Find the Go module name and path
	goModuleName, goModulePath, err := module.FindGoModuleNameAndPath(currentDir)
	if err != nil {
		log.Fatalf("error occurred finding module path: %s", err.Error())
		os.Exit(1)
	}

	// Load the configuration from the gio.yaml file
	conf, err := config.Load("gio.yaml")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Build the Regular Expressions for excluding files/folders
	excludeByNameRegExp, excludeByPathRegExp := excludes.Build(conf.Excludes)

	// Build the Regular Expressions and DisplayOrder for the group definitions
	groupRegExpMatchers, displayOrder := groups.Build(conf.Groups, goModuleName)

	wg.Add(1)
	// Start up the Format worker so that it is ready when we start queuing up files
	go imports.Format(&files, &wg, groupRegExpMatchers, displayOrder, listOnly)

	// Set the basePath for use later
	basePath := goModulePath + "/"

	// Change our working directory to the goModulePath
	err = os.Chdir(goModulePath)
	if err != nil {
		panic(err)
	}

	// Pre-optization so that we can skip the Name or Path matches if they are empty
	excludeByNameRegExpLenOk := len(excludeByNameRegExp.String()) != 0
	excludeByPathRegExpLenOk := len(excludeByPathRegExp.String()) != 0

	// Walk the filepath looking for Go files
	if err = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := f.Name()
		isDir := f.IsDir()
		isGoFile := strings.HasSuffix(name, ".go")
		relativePath := strings.Replace(path, basePath, "", 1)
		// If the object is not a directory and not a Go file, skip it
		if isDir || isGoFile {
			// If the objects name or path matches an exclude Regular Expression, skip it
			if (excludeByNameRegExpLenOk && excludeByNameRegExp.MatchString(name)) || (excludeByPathRegExpLenOk && excludeByPathRegExp.MatchString(relativePath)) {
				if isDir {
					return filepath.SkipDir
				}
				return nil
			}
			// If the object is a Go file and is not excluded, queue it for formatting
			if isGoFile {
				files <- relativePath
			}
		}
		return nil
	}); err != nil {
		log.Fatalf("unable to complete walking file tree: %s", err.Error())
	}

	// Close the files channel since we are done queuing up files to format
	close(files)

	// Wait for all files to be processed
	wg.Wait()

	// set MEMPROFILE=<filename> to create a <filename>.pprof memory profile file
	if len(os.Getenv("MEMPROFILE")) != 0 {
		f, err := os.Create(os.Getenv("MEMPROFILE"))
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
