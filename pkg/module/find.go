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
package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func FindGoModuleAndPath(path string) (string, string, error) {
	if s, err := os.Stat(path); err != nil {
		return "", "", err
	} else if !s.IsDir() {
		path = filepath.Dir(path)
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}

	for path != "." {
		modFilePath := fmt.Sprintf("%s/go.mod", path)
		if _, err := os.Stat(modFilePath); !os.IsNotExist(err) {
			break
		}
		prevPath := path
		path = filepath.Dir(path)
		if path == prevPath {
			return "", "", nil
		}
	}
	if path == "." {
		return "", "", nil
	}

	f, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", path))
	if err != nil {
		return "", "", fmt.Errorf("unable to open go.mod file for reading: %v", err)
	}
	module := modfile.ModulePath(f)
	if len(module) == 0 {
		return "", "", fmt.Errorf("unable to automatically determine module path, please provide one using the --module flag")
	}
	return module, path, nil
}
