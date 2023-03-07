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
	"strings"
	"testing"
)

func TestFindGoModuleAndPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		wantModule string
		wantPath   string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "path does not exist",
			args: args{
				path: "does/not/exist",
			},
			wantModule: "",
			wantPath:   "",
			wantErr:    true,
		},
		{
			name: "path is a file not a directory",
			args: args{
				path: "../../test/testdata/findModule/moduleOne/main.go",
			},
			wantModule: "github.com/example/moduleOne",
			wantPath:   "test/testdata/findModule/moduleOne",
			wantErr:    false,
		},
		{
			name: "path is a subdirectory",
			args: args{
				path: "../../test/testdata/findModule/moduleOne/folderOne/folderTwo",
			},
			wantModule: "github.com/example/moduleOne",
			wantPath:   "test/testdata/findModule/moduleOne",
			wantErr:    false,
		},
		{
			name: "path does not exist",
			args: args{
				path: "../../test/testdata/findModule/moduleFour",
			},
			wantModule: "",
			wantPath:   "",
			wantErr:    true,
			wantErrMsg: "moduleFour does not exist",
		},
		{
			name: "unable to determine path",
			args: args{
				path: "../../test/testdata/findModule/moduleThree",
			},
			wantModule: "",
			wantPath:   "",
			wantErr:    true,
			wantErrMsg: "unable to determine module",
		},
		{
			name: "path is module base directory",
			args: args{
				path: "../../test/testdata/findModule/moduleOne",
			},
			wantModule: "github.com/example/moduleOne",
			wantPath:   "test/testdata/findModule/moduleOne",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotModule, gotPath, gotErr := FindGoModuleNameAndPath(tt.args.path)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("FindGoModuleAndPath() error = %v, wantErr = %v", gotErr, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !strings.Contains(gotErr.Error(), tt.wantErrMsg) {
					t.Errorf("FindGoModuleAndPath() gotErrMsg = %v, wantErrMsg = %v", gotErr.Error(), tt.wantErrMsg)
				}
			}
			if gotModule != tt.wantModule {
				t.Errorf("FindGoModuleAndPath() gotModule = %v, wantModule = %v", gotModule, tt.wantModule)
			}
			if !strings.Contains(gotPath, tt.wantPath) {
				t.Errorf("FindGoModuleAndPath() gotPath = %v, wantPath = %v", gotPath, tt.wantPath)
			}
		})
	}
}
