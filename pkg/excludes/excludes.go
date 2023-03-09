package excludes

import (
	"regexp"
	"strings"

	v1alpha1 "github.com/go-imports-organizer/goio/pkg/api/v1alpha1"
)

func Build(excludes []v1alpha1.Exclude) (*regexp.Regexp, *regexp.Regexp) {
	var excludeByPath []string
	var excludeByName []string

	for _, exclude := range excludes {
		switch exclude.MatchType {
		case v1alpha1.ExcludeMatchTypeName:
			excludeByName = append(excludeByName, exclude.RegExp)
		case v1alpha1.ExcludeMatchTypeRelativePath:
			excludeByPath = append(excludeByPath, exclude.RegExp)
		}
	}
	return regexp.MustCompile(strings.Join(excludeByName, "|")), regexp.MustCompile(strings.Join(excludeByPath, "|"))
}
