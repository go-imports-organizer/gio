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
	"bytes"
	"go/ast"
	"reflect"
	"sync"
	"testing"

	v1alpha1 "github.com/go-imports-organizer/goio/pkg/api/v1alpha1"
	"github.com/go-imports-organizer/goio/pkg/groups"
)

func TestAddSpaces(t *testing.T) {
	type args struct {
		input  []byte
		breaks []string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "works with function",
			args: args{
				input: []byte(`imports (
	"fmt",
	"io",
	"reflect",
	"sort",
)

func main() {}
`),
				breaks: []string{
					"io",
					"sort",
				},
			},
			want: []byte(`imports (
	"fmt",

	"io",
	"reflect",

	"sort",
)

func main() {}
`),
		},
		{
			name: "standard imports",
			args: args{
				input: []byte(`imports (
	"fmt",
	"io",
	"reflect",
	"sort",
)
`),
				breaks: []string{
					"io",
					"sort",
				},
			},
			want: []byte(`imports (
	"fmt",

	"io",
	"reflect",

	"sort",
)
`),
		},
		{
			name: "kubernetes imports",
			args: args{
				input: []byte(`imports (
	"k8s.io/api/core/v1alpha1",
	"k8s.io/apimachinery/pkg/api/errors",
	"k8s.io/apimachinery/pkg/apis/meta/v1alpha1",
	"k8s.io/apimachinery/pkg/util/sets",
	"k8s.io/apimachinery/pkg/util/uuid",
	"k8s.io/apimachinery/pkg/util/wait",
)
`),
				breaks: []string{
					"k8s.io/apimachinery/pkg/apis/meta/v1alpha1",
					"k8s.io/apimachinery/pkg/util/uuid",
				},
			},
			want: []byte(`imports (
	"k8s.io/api/core/v1alpha1",
	"k8s.io/apimachinery/pkg/api/errors",

	"k8s.io/apimachinery/pkg/apis/meta/v1alpha1",
	"k8s.io/apimachinery/pkg/util/sets",

	"k8s.io/apimachinery/pkg/util/uuid",
	"k8s.io/apimachinery/pkg/util/wait",
)
`),
		},
		{
			name: "openshift imports",
			args: args{
				input: []byte(`imports (
	"github.com/openshift/api/build/v1alpha1",
	"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1alpha1",
	"github.com/openshift/imagebuilder",
	"github.com/openshift/imagebuilder/dockerfile/command",
	"github.com/openshift/imagebuilder/dockerfile/parser",
	"github.com/openshift/library-go/pkg/git",
	"github.com/openshift/library-go/pkg/image/reference",
	"github.com/openshift/source-to-image/pkg/scm/git",
	"github.com/openshift/source-to-image/pkg/util",
)
`),
				breaks: []string{
					"github.com/openshift/imagebuilder/dockerfile/command",
					"github.com/openshift/library-go/pkg/git",
					"github.com/openshift/source-to-image/pkg/util",
				},
			},
			want: []byte(`imports (
	"github.com/openshift/api/build/v1alpha1",
	"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1alpha1",
	"github.com/openshift/imagebuilder",

	"github.com/openshift/imagebuilder/dockerfile/command",
	"github.com/openshift/imagebuilder/dockerfile/parser",

	"github.com/openshift/library-go/pkg/git",
	"github.com/openshift/library-go/pkg/image/reference",
	"github.com/openshift/source-to-image/pkg/scm/git",

	"github.com/openshift/source-to-image/pkg/util",
)
`),
		},
		{
			name: "mixed imports",
			args: args{
				input: []byte(`imports (
	"fmt",
	"io",
	"reflect",
	"sort",
	"k8s.io/api/core/v1alpha1",
	"k8s.io/apimachinery/pkg/api/errors",
	"k8s.io/apimachinery/pkg/apis/meta/v1alpha1",
	"k8s.io/apimachinery/pkg/util/sets",
	"k8s.io/apimachinery/pkg/util/uuid",
	"k8s.io/apimachinery/pkg/util/wait",
	"github.com/openshift/api/build/v1alpha1",
	"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1alpha1",
	"github.com/openshift/imagebuilder",
	"github.com/openshift/imagebuilder/dockerfile/command",
	"github.com/openshift/imagebuilder/dockerfile/parser",
	"github.com/openshift/library-go/pkg/git",
	"github.com/openshift/library-go/pkg/image/reference",
	"github.com/openshift/source-to-image/pkg/scm/git",
	"github.com/openshift/source-to-image/pkg/util",
)
`),
				breaks: []string{
					"k8s.io/apimachinery/pkg/util/sets",
					"github.com/openshift/api/build/v1alpha1",
					"github.com/openshift/imagebuilder/dockerfile/parser",
					"github.com/openshift/source-to-image/pkg/scm/git",
				},
			},
			want: []byte(`imports (
	"fmt",
	"io",
	"reflect",
	"sort",
	"k8s.io/api/core/v1alpha1",
	"k8s.io/apimachinery/pkg/api/errors",
	"k8s.io/apimachinery/pkg/apis/meta/v1alpha1",

	"k8s.io/apimachinery/pkg/util/sets",
	"k8s.io/apimachinery/pkg/util/uuid",
	"k8s.io/apimachinery/pkg/util/wait",

	"github.com/openshift/api/build/v1alpha1",
	"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1alpha1",
	"github.com/openshift/imagebuilder",
	"github.com/openshift/imagebuilder/dockerfile/command",

	"github.com/openshift/imagebuilder/dockerfile/parser",
	"github.com/openshift/library-go/pkg/git",
	"github.com/openshift/library-go/pkg/image/reference",

	"github.com/openshift/source-to-image/pkg/scm/git",
	"github.com/openshift/source-to-image/pkg/util",
)
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			buf.Write(tt.args.input)
			got, err := AddSpaces(bytes.NewReader(buf.Bytes()), tt.args.breaks)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("AddSpaces() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type args struct {
		regExpMatchers []v1alpha1.RegExpMatcher
		displayOrder   []string
		listOnly       *bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			files := make(chan string, 10)
			resultsChan := make(chan string)
			hasResults := false
			defer close(files)
			Format(&files, &resultsChan, &hasResults, &wg, tt.args.regExpMatchers, tt.args.displayOrder, tt.args.listOnly)
			wg.Wait()
		})
	}
}

func TestPopulateGroups(t *testing.T) {
	type args struct {
		imports      []*ast.ImportSpec
		groups       []v1alpha1.Group
		goModuleName string
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantImportGroups map[string][]ast.ImportSpec
	}{
		{
			name: "testing",
			args: args{
				imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Value: ``,
						},
					},
					{
						Path: &ast.BasicLit{
							Value: `"github.com/exampleOne/module/pkg/packageOne"`,
						},
					},
					{
						Path: &ast.BasicLit{
							Value: `"github.com/exampleTwo/module/pkg/packageOne"`,
						},
					},
				},
				groups: []v1alpha1.Group{
					{
						MatchOrder:  0,
						Description: "module",
						RegExp:      []string{"%{module}%"},
					},
					{
						MatchOrder:  1,
						Description: "standard",
						RegExp:      []string{`^[a-zA-Z0-9\\/]+$`},
					},
					{
						MatchOrder:  2,
						Description: "other",
						RegExp:      []string{`[a-zA-Z0-9]+\\.[a-zA-Z0-9]+/`},
					},
				},
				goModuleName: "github.com/exampleOne/module",
			},
			wantErr: false,
			wantImportGroups: map[string][]ast.ImportSpec{
				"module": {
					{
						Path: &ast.BasicLit{
							Value: `"github.com/exampleOne/module/pkg/packageOne"`,
						},
					},
				},
				"other": {
					{
						Path: &ast.BasicLit{
							Value: `"github.com/exampleTwo/module/pkg/packageOne"`,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		importGroups := make(map[string][]ast.ImportSpec)
		groupRegExpMatchers, _ := groups.Build(tt.args.groups, tt.args.goModuleName)
		t.Run(tt.name, func(t *testing.T) {
			if err := PopulateGroups(importGroups, groupRegExpMatchers, tt.args.imports); (err != nil) != tt.wantErr {
				t.Errorf("PopulateGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
			for wantGroup, wantImports := range tt.wantImportGroups {
				if gotImports, ok := importGroups[wantGroup]; ok {
					for _, i := range wantImports {
						found := false
						for _, j := range gotImports {
							if i.Path.Value == j.Path.Value {
								found = true
							}
						}
						if !found {
							t.Errorf("%#v not found in %#v", i.Path.Value, gotImports)
						}
					}
				} else {
					t.Errorf("group %s not found in importGroups: %#v", wantGroup, importGroups)
				}
			}
		})
	}
}

func TestInsertGroups(t *testing.T) {
	type args struct {
		f            *ast.File
		importGroups map[string][]ast.ImportSpec
		displayOrder []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertGroups(tt.args.f, tt.args.importGroups, tt.args.displayOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}
