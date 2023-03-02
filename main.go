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
	"regexp"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	v1 "github.com/go-imports-organizer/gio/pkg/api/v1"
	"github.com/go-imports-organizer/gio/pkg/imports"
	"github.com/go-imports-organizer/gio/pkg/module"
	"github.com/go-imports-organizer/gio/pkg/sorter"
)

var (
	wg    sync.WaitGroup
	files = make(chan string)
)

func main() {
	listOnly := flag.Bool("l", false, "list files that need to be organized (no changes made)")
	flag.Parse()

	if len(os.Getenv("CPUPROFILE")) != 0 {
		f, err := os.Create(os.Getenv("CPUPROFILE"))
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	goModuleName, goModulePath, err := module.FindGoModuleAndPath(".")
	if err != nil {
		log.Fatalf("error occurred finding module path: %s", err.Error())
		os.Exit(1)
	}

	var configFile []byte
	if configFile, err = os.ReadFile("gio.yaml"); err != nil {
		log.Fatalf("unable to read configuration file gio.yaml: %s", err.Error())
	}

	var config v1.Config
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("unable to unmarshal file gio.yaml: %s", err.Error())
	}

	var excludeByPath []string
	var excludeByName []string

	for _, exclude := range config.Excludes {
		switch exclude.MatchType {
		case v1.ExcludeMatchTypeName:
			excludeByName = append(excludeByName, exclude.RegExp)
		case v1.ExcludeMatchTypeRelativePath:
			excludeByPath = append(excludeByPath, exclude.RegExp)
		}
	}
	excludeByNameRegExp := regexp.MustCompile(strings.Join(excludeByName, "|"))
	excludeByPathRegExp := regexp.MustCompile(strings.Join(excludeByPath, "|"))

	groupRegExpMatchers := []v1.RegExpMatcher{}
	displayOrder := []string{}

	sort.Sort(sorter.SortGroupsByDisplayOrder(config.Groups))
	for _, group := range config.Groups {
		displayOrder = append(displayOrder, group.Description)
	}

	sort.Sort(sorter.SortGroupsByMatchOrder(config.Groups))

	for i := range config.Groups {
		if config.Groups[i].RegExp == `%{module}%` {
			config.Groups[i].RegExp = fmt.Sprintf("^%s", strings.ReplaceAll(strings.ReplaceAll(goModuleName, `.`, `\.`), `/`, `\/`))
		}
		groupRegExpMatchers = append(groupRegExpMatchers, v1.RegExpMatcher{
			Bucket: config.Groups[i].Description,
			RegExp: regexp.MustCompile(config.Groups[i].RegExp),
		},
		)
	}

	wg.Add(1)
	go imports.Format(&files, &wg, groupRegExpMatchers, displayOrder, listOnly)

	basePath := goModulePath + "/"

	err = os.Chdir(goModulePath)
	if err != nil {
		panic(err)
	}

	excludeByNameRegExpLenOk := len(excludeByNameRegExp.String()) != 0
	excludeByPathRegExpLenOk := len(excludeByPathRegExp.String()) != 0

	if err = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := f.Name()
		isDir := f.IsDir()
		isGoFile := strings.HasSuffix(name, ".go")
		relativePath := strings.Replace(path, basePath, "", 1)
		if isDir || isGoFile {
			if (excludeByNameRegExpLenOk && excludeByNameRegExp.MatchString(name)) || (excludeByPathRegExpLenOk && excludeByPathRegExp.MatchString(relativePath)) {
				if isDir {
					return filepath.SkipDir
				}
				return nil
			}
			if isGoFile {
				files <- relativePath
			}
		}
		return nil
	}); err != nil {
		log.Fatalf("unable to complete walking file tree: %s", err.Error())
	}

	close(files)

	wg.Wait()

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
