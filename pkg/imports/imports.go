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
package imports

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	v1 "github.com/go-imports-organizer/gio/pkg/api/v1"
	"github.com/go-imports-organizer/gio/pkg/sorter"
)

// taken from https://github.com/golang/tools/blob/71482053b885ea3938876d1306ad8a1e4037f367/internal/imports/imports.go#L380
func addSpaces(r io.Reader, breaks []string) ([]byte, error) {
	var out bytes.Buffer
	in := bufio.NewReader(r)
	inImports := false
	done := false
	for {
		s, err := in.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !inImports && !done && strings.HasPrefix(s, "import") {
			inImports = true
		}
		if inImports && (strings.HasPrefix(s, "var") ||
			strings.HasPrefix(s, "func") ||
			strings.HasPrefix(s, "const") ||
			strings.HasPrefix(s, "type")) {
			done = true
			inImports = false
		}
		if inImports && len(breaks) > 0 {
			if m := regexp.MustCompile(`^\s+(?:[\w\.]+\s+)?"(.+)"`).FindStringSubmatch(s); m != nil {
				if m[1] == breaks[0] {
					out.WriteByte('\n')
					breaks = breaks[1:]
				}
			}
		}

		fmt.Fprint(&out, s)
	}
	return out.Bytes(), nil
}

// Format takes a channel of file paths and formats the files imports
func Format(files *chan string, wg *sync.WaitGroup, regExpMatchers []v1.RegExpMatcher, displayOrder []string, listOnly *bool) {
	defer wg.Done()
	for path := range *files {
		var importGroups = make(map[string][]ast.ImportSpec)
		if len(path) == 0 {
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		oldModTime := info.ModTime()

		var breaks []string

		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
		if err != nil {
			var scannerErrorList scanner.ErrorList
			if errors.As(err, &scannerErrorList) {
				for _, err := range scannerErrorList {
					log.Fatalf("%s", err)
				}
			} else {
				log.Fatalf("%s", err.Error())
			}
		}
		for _, i := range f.Imports {
			if len(i.Path.Value) == 0 {
				continue
			}
			found := false
			unquotedPath, err := strconv.Unquote(i.Path.Value)
			if err != nil {
				log.Printf("unable to unquote %s", i.Path.Value)
			}
			for _, r := range regExpMatchers {
				if r.RegExp.MatchString(unquotedPath) {
					if _, ok := importGroups[r.Bucket]; !ok {
						importGroups[r.Bucket] = []ast.ImportSpec{}
					}
					importGroups[r.Bucket] = append(importGroups[r.Bucket], *i)
					found = true
					break
				}
			}
			if !found {
				importGroups["other"] = append(importGroups["other"], *i)
			}
		}
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if ok && gen.Tok == token.IMPORT {
				gen.Specs = []ast.Spec{}
				for _, group := range displayOrder {
					sort.Sort(sorter.ByPathValue(importGroups[group]))
					for n := range importGroups[group] {
						importGroups[group][n].EndPos = 0
						importGroups[group][n].Path.ValuePos = 0
						if importGroups[group][n].Name != nil {
							importGroups[group][n].Name.NamePos = 0
						}
						gen.Specs = append(gen.Specs, &importGroups[group][n])
						if n == 0 && group != displayOrder[0] {
							newstr, err := strconv.Unquote(importGroups[group][n].Path.Value)
							if err != nil {
								log.Fatalf("%#v", err)
							}
							breaks = append(breaks, newstr)
						}
					}
				}
			}
		}

		printerMode := printer.TabIndent

		printConfig := &printer.Config{Mode: printerMode, Tabwidth: 4}

		var buf bytes.Buffer
		if err = printConfig.Fprint(&buf, fs, f); err != nil {
			log.Fatalf("%s", err.Error())
		}
		out, err := addSpaces(bytes.NewReader(buf.Bytes()), breaks)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		out, err = format.Source(out)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		if *listOnly {
			oldFile, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("unable to read file %q: %s", path, err.Error())
			}
			if !bytes.Equal(oldFile, out) {
				log.Printf("%s is not sorted \n", path)
			}
		}

		info, err = os.Stat(path)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		if !info.ModTime().Equal(oldModTime) {
			log.Printf("%s was modified while formatting, cowardly refusing to overwrite", path)
			continue
		}
		if err = ioutil.WriteFile(path, out, info.Mode()); err != nil {
			log.Fatalf("%#v", err)
		}

	}

}
