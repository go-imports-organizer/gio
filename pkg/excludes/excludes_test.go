package excludes

import (
	"reflect"
	"regexp"
	"testing"

	v1 "github.com/go-imports-organizer/goio/pkg/api/v1"
)

func TestBuild(t *testing.T) {
	type args struct {
		excludes []v1.Exclude
	}
	tests := []struct {
		name             string
		args             args
		wantNameMatchers *regexp.Regexp
		wantPathMatchers *regexp.Regexp
	}{
		{
			name: "only name excludes",
			args: args{
				excludes: []v1.Exclude{
					{
						MatchType: "name",
						RegExp:    "^name-one$",
					},
					{
						MatchType: "name",
						RegExp:    "^name-two$",
					},
				},
			},
			wantNameMatchers: regexp.MustCompile(`^name-one$|^name-two$`),
			wantPathMatchers: regexp.MustCompile(``),
		},
		{
			name: "only path excludes",
			args: args{
				excludes: []v1.Exclude{
					{
						MatchType: "path",
						RegExp:    "^path-one$",
					},
					{
						MatchType: "path",
						RegExp:    "^path-two$",
					},
				},
			},
			wantNameMatchers: regexp.MustCompile(``),
			wantPathMatchers: regexp.MustCompile(`^path-one$|^path-two$`),
		},
		{
			name: "name and path excludes",
			args: args{
				excludes: []v1.Exclude{
					{
						MatchType: "name",
						RegExp:    "^name-one$",
					},
					{
						MatchType: "name",
						RegExp:    "^name-two$",
					},
					{
						MatchType: "path",
						RegExp:    "^path-one$",
					},
					{
						MatchType: "path",
						RegExp:    "^path-two$",
					},
				},
			},
			wantNameMatchers: regexp.MustCompile(`^name-one$|^name-two$`),
			wantPathMatchers: regexp.MustCompile(`^path-one$|^path-two$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNameMatcher, gotPathMatcher := Build(tt.args.excludes)
			if !reflect.DeepEqual(gotNameMatcher, tt.wantNameMatchers) {
				t.Errorf("Build() gotNameMatcher = %v, want %v", gotNameMatcher, tt.wantNameMatchers)
			}
			if !reflect.DeepEqual(gotPathMatcher, tt.wantPathMatchers) {
				t.Errorf("Build() gotPathMatcher = %v, want %v", gotPathMatcher, tt.wantPathMatchers)
			}
		})
	}
}
