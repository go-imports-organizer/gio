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
package sorter

import (
	"go/ast"
	"reflect"
	"sort"
	"testing"

	v1 "github.com/go-imports-organizer/gio/pkg/api/v1"
)

func TestSortImportsByPathValue(t *testing.T) {
	type args struct {
		imports []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "sort standard imports",
			args: args{
				imports: []string{
					"sort",
					"reflect",
					"fmt",
					"io",
				},
			},
			want: []string{
				"fmt",
				"io",
				"reflect",
				"sort",
			},
		},
		{
			name: "sort kubernetes imports",
			args: args{
				imports: []string{
					"k8s.io/apimachinery/pkg/api/errors",
					"k8s.io/apimachinery/pkg/util/wait",
					"k8s.io/apimachinery/pkg/apis/meta/v1",
					"k8s.io/apimachinery/pkg/util/uuid",
					"k8s.io/apimachinery/pkg/util/sets",
					"k8s.io/api/core/v1",
				},
			},
			want: []string{
				"k8s.io/api/core/v1",
				"k8s.io/apimachinery/pkg/api/errors",
				"k8s.io/apimachinery/pkg/apis/meta/v1",
				"k8s.io/apimachinery/pkg/util/sets",
				"k8s.io/apimachinery/pkg/util/uuid",
				"k8s.io/apimachinery/pkg/util/wait",
			},
		},
		{
			name: "sort openshift imports",
			args: args{
				imports: []string{
					"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1",
					"github.com/openshift/library-go/pkg/git",
					"github.com/openshift/imagebuilder",
					"github.com/openshift/source-to-image/pkg/util",
					"github.com/openshift/imagebuilder/dockerfile/command",
					"github.com/openshift/api/build/v1",
					"github.com/openshift/library-go/pkg/image/reference",
					"github.com/openshift/imagebuilder/dockerfile/parser",
					"github.com/openshift/source-to-image/pkg/scm/git",
				},
			},
			want: []string{
				"github.com/openshift/api/build/v1",
				"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1",
				"github.com/openshift/imagebuilder",
				"github.com/openshift/imagebuilder/dockerfile/command",
				"github.com/openshift/imagebuilder/dockerfile/parser",
				"github.com/openshift/library-go/pkg/git",
				"github.com/openshift/library-go/pkg/image/reference",
				"github.com/openshift/source-to-image/pkg/scm/git",
				"github.com/openshift/source-to-image/pkg/util",
			},
		},
		{
			name: "sort mixed imports",
			args: args{
				imports: []string{
					"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1",
					"github.com/openshift/library-go/pkg/git",
					"github.com/openshift/imagebuilder",
					"github.com/openshift/source-to-image/pkg/util",
					"github.com/openshift/imagebuilder/dockerfile/command",
					"github.com/openshift/api/build/v1",
					"github.com/openshift/library-go/pkg/image/reference",
					"github.com/openshift/imagebuilder/dockerfile/parser",
					"github.com/openshift/source-to-image/pkg/scm/git",
					"sort",
					"reflect",
					"fmt",
					"io",
					"k8s.io/apimachinery/pkg/api/errors",
					"k8s.io/apimachinery/pkg/util/wait",
					"k8s.io/apimachinery/pkg/apis/meta/v1",
					"k8s.io/apimachinery/pkg/util/uuid",
					"k8s.io/apimachinery/pkg/util/sets",
					"k8s.io/api/core/v1",
				},
			},
			want: []string{
				"fmt",
				"github.com/openshift/api/build/v1",
				"github.com/openshift/client-go/build/clientset/versioned/typed/build/v1",
				"github.com/openshift/imagebuilder",
				"github.com/openshift/imagebuilder/dockerfile/command",
				"github.com/openshift/imagebuilder/dockerfile/parser",
				"github.com/openshift/library-go/pkg/git",
				"github.com/openshift/library-go/pkg/image/reference",
				"github.com/openshift/source-to-image/pkg/scm/git",
				"github.com/openshift/source-to-image/pkg/util",
				"io",
				"k8s.io/api/core/v1",
				"k8s.io/apimachinery/pkg/api/errors",
				"k8s.io/apimachinery/pkg/apis/meta/v1",
				"k8s.io/apimachinery/pkg/util/sets",
				"k8s.io/apimachinery/pkg/util/uuid",
				"k8s.io/apimachinery/pkg/util/wait",
				"reflect",
				"sort",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imports := []ast.ImportSpec{}
			for _, s := range tt.args.imports {
				is := ast.ImportSpec{
					Path: &ast.BasicLit{
						Value: s,
					},
				}
				imports = append(imports, is)
			}
			sort.Sort(SortImportsByPathValue(imports))
			sortedImports := []string{}
			for _, s := range imports {
				sortedImports = append(sortedImports, s.Path.Value)
			}
			if !reflect.DeepEqual(sortedImports, tt.want) {
				t.Errorf("SortImportsByPathValue() = %v, want %v", sortedImports, tt.want)
			}
		})
	}
}

func TestSortGroupsByMatchOrder(t *testing.T) {
	type args struct {
		groups []v1.Group
	}
	tests := []struct {
		name string
		args args
		want []v1.Group
	}{
		{
			name: "sort groups one",
			args: args{
				groups: []v1.Group{
					{
						MatchOrder:   3,
						DisplayOrder: 3,
						Description:  "Group 3",
						RegExp:       "group-3-regexp",
					},
					{
						MatchOrder:   5,
						DisplayOrder: 1,
						Description:  "Group 5",
						RegExp:       "group-5-regexp",
					},
					{
						MatchOrder:   2,
						DisplayOrder: 4,
						Description:  "Group 2",
						RegExp:       "group-2-regexp",
					},
					{
						MatchOrder:   1,
						DisplayOrder: 5,
						Description:  "Group 1",
						RegExp:       "group-1-regexp",
					},
					{
						MatchOrder:   4,
						DisplayOrder: 2,
						Description:  "Group 4",
						RegExp:       "group-4-regexp",
					},
				},
			},
			want: []v1.Group{
				{
					MatchOrder:   1,
					DisplayOrder: 5,
					Description:  "Group 1",
					RegExp:       "group-1-regexp",
				},
				{
					MatchOrder:   2,
					DisplayOrder: 4,
					Description:  "Group 2",
					RegExp:       "group-2-regexp",
				},
				{
					MatchOrder:   3,
					DisplayOrder: 3,
					Description:  "Group 3",
					RegExp:       "group-3-regexp",
				},
				{
					MatchOrder:   4,
					DisplayOrder: 2,
					Description:  "Group 4",
					RegExp:       "group-4-regexp",
				},
				{
					MatchOrder:   5,
					DisplayOrder: 1,
					Description:  "Group 5",
					RegExp:       "group-5-regexp",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(SortGroupsByMatchOrder(tt.args.groups))

			if !reflect.DeepEqual(tt.args.groups, tt.want) {
				t.Errorf("SortGroupsByMatchOrder() = %v, want %v", tt.args.groups, tt.want)
			}
		})
	}
}

func TestSortGroupsByDisplayOrder(t *testing.T) {
	type args struct {
		groups []v1.Group
	}
	tests := []struct {
		name string
		args args
		want []v1.Group
	}{
		{
			name: "sort groups one",
			args: args{
				groups: []v1.Group{
					{
						MatchOrder:   3,
						DisplayOrder: 3,
						Description:  "Group 3",
						RegExp:       "group-3-regexp",
					},
					{
						MatchOrder:   5,
						DisplayOrder: 1,
						Description:  "Group 5",
						RegExp:       "group-5-regexp",
					},
					{
						MatchOrder:   2,
						DisplayOrder: 4,
						Description:  "Group 2",
						RegExp:       "group-2-regexp",
					},
					{
						MatchOrder:   1,
						DisplayOrder: 5,
						Description:  "Group 1",
						RegExp:       "group-1-regexp",
					},
					{
						MatchOrder:   4,
						DisplayOrder: 2,
						Description:  "Group 4",
						RegExp:       "group-4-regexp",
					},
				},
			},
			want: []v1.Group{
				{
					MatchOrder:   5,
					DisplayOrder: 1,
					Description:  "Group 5",
					RegExp:       "group-5-regexp",
				},
				{
					MatchOrder:   4,
					DisplayOrder: 2,
					Description:  "Group 4",
					RegExp:       "group-4-regexp",
				},
				{
					MatchOrder:   3,
					DisplayOrder: 3,
					Description:  "Group 3",
					RegExp:       "group-3-regexp",
				},
				{
					MatchOrder:   2,
					DisplayOrder: 4,
					Description:  "Group 2",
					RegExp:       "group-2-regexp",
				},
				{
					MatchOrder:   1,
					DisplayOrder: 5,
					Description:  "Group 1",
					RegExp:       "group-1-regexp",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(SortGroupsByDisplayOrder(tt.args.groups))

			if !reflect.DeepEqual(tt.args.groups, tt.want) {
				t.Errorf("SortGroupsByDisplayOrder() = %v, want %v", tt.args.groups, tt.want)
			}
		})
	}
}
