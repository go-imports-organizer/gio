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
package v1

import (
	"regexp"
)

type RegExpMatcher struct {
	Bucket string
	RegExp *regexp.Regexp
}

const (
	ExcludeMatchTypeName         string = "name"
	ExcludeMatchTypeRelativePath string = "path"
)

type Exclude struct {
	MatchType string
	RegExp    string
}

type Group struct {
	MatchOrder   int
	DisplayOrder int
	Description  string
	RegExp       string
}

const (
	GroupMatchValueStandard          string = "standard"
	GroupMatchValueModule            string = "module"
	GroupMatchValueOther             string = "other"
	GroupMatchValueModulePlaceholder string = "%{module}%"
)

type Config struct {
	Excludes []Exclude
	Groups   []Group
}
