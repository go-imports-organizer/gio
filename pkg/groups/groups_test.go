package groups

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	v1alpha1 "github.com/go-imports-organizer/goio/pkg/api/v1alpha1"
)

func TestBuild(t *testing.T) {
	type args struct {
		groups       []v1alpha1.Group
		goModuleName string
	}
	tests := []struct {
		name               string
		args               args
		wantRegExpMatchers []v1alpha1.RegExpMatcher
		wantDisplayOrder   []string
	}{
		{
			name: "group one test",
			args: args{
				goModuleName: "github.com/example/module",
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
			},
			wantRegExpMatchers: []v1alpha1.RegExpMatcher{
				{
					Bucket: "module",
					RegExp: regexp.MustCompile(fmt.Sprintf("^%s", strings.ReplaceAll(strings.ReplaceAll(`github.com/example/module`, `.`, `\.`), `/`, `\/`))),
				},
				{
					Bucket: "standard",
					RegExp: regexp.MustCompile(`^[a-zA-Z0-9\\/]+$`),
				},
				{
					Bucket: "other",
					RegExp: regexp.MustCompile(`[a-zA-Z0-9]+\\.[a-zA-Z0-9]+/`),
				},
			},
			wantDisplayOrder: []string{
				"standard",
				"other",
				"module",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRegExpMatchers, _ := Build(tt.args.groups, tt.args.goModuleName)
			if !reflect.DeepEqual(gotRegExpMatchers, tt.wantRegExpMatchers) {
				t.Errorf("Build() gotRegExpMatchers = %v, wantRegExpMatchers %v", gotRegExpMatchers, tt.wantRegExpMatchers)
			}
		})
	}
}
