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

	v1 "github.com/go-imports-organizer/gio/pkg/api/v1"
)

type SortImportsByPathValue []ast.ImportSpec

func (a SortImportsByPathValue) Len() int           { return len(a) }
func (a SortImportsByPathValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortImportsByPathValue) Less(i, j int) bool { return a[i].Path.Value < a[j].Path.Value }

type SortGroupsByMatchOrder []v1.Group

func (a SortGroupsByMatchOrder) Len() int           { return len(a) }
func (a SortGroupsByMatchOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGroupsByMatchOrder) Less(i, j int) bool { return a[i].MatchOrder < a[j].MatchOrder }

type SortGroupsByDisplayOrder []v1.Group

func (a SortGroupsByDisplayOrder) Len() int           { return len(a) }
func (a SortGroupsByDisplayOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGroupsByDisplayOrder) Less(i, j int) bool { return a[i].DisplayOrder < a[j].DisplayOrder }
