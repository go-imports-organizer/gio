package groups

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	v1 "github.com/go-imports-organizer/goio/pkg/api/v1"
	"github.com/go-imports-organizer/goio/pkg/sorter"
)

func Build(groups []v1.Group, goModuleName string) ([]v1.RegExpMatcher, []string) {
	groupRegExpMatchers := []v1.RegExpMatcher{}
	displayOrder := []string{}

	sort.Sort(sorter.SortGroupsByDisplayOrder(groups))
	for _, group := range groups {
		displayOrder = append(displayOrder, group.Description)
	}

	sort.Sort(sorter.SortGroupsByMatchOrder(groups))

	for i := range groups {
		if groups[i].RegExp == `%{module}%` {
			groups[i].RegExp = fmt.Sprintf("^%s", strings.ReplaceAll(strings.ReplaceAll(goModuleName, `.`, `\.`), `/`, `\/`))
		}
		groupRegExpMatchers = append(groupRegExpMatchers, v1.RegExpMatcher{
			Bucket: groups[i].Description,
			RegExp: regexp.MustCompile(groups[i].RegExp),
		},
		)
	}
	return groupRegExpMatchers, displayOrder
}
