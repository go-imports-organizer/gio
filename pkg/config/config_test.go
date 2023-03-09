package config

import (
	"reflect"
	"strings"
	"testing"

	v1 "github.com/go-imports-organizer/goio/pkg/api/v1"
)

func TestLoad(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name       string
		args       args
		want       v1.Config
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "working config file",
			args: args{
				file: "../../test/testdata/config/works.yaml",
			},
			want: v1.Config{
				Excludes: []v1.Exclude{
					{
						MatchType: "name",
						RegExp:    "^\\.git$",
					},
					{
						MatchType: "name",
						RegExp:    "^vendor$",
					},
				},
				Groups: []v1.Group{
					{
						MatchOrder:   0,
						DisplayOrder: 2,
						Description:  "module",
						RegExp:       "%{module}%",
					},
					{
						MatchOrder:   1,
						DisplayOrder: 0,
						Description:  "standard",
						RegExp:       "^[a-zA-Z0-9\\/]+$",
					},
					{
						MatchOrder:   2,
						DisplayOrder: 1,
						Description:  "other",
						RegExp:       "[a-zA-Z0-9]+\\.[a-zA-Z0-9]+/",
					},
				},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "malformed yaml file",
			args: args{
				file: "../../test/testdata/config/malformed.yaml",
			},
			want:       v1.Config{},
			wantErr:    true,
			wantErrMsg: "unable to unmarshal file",
		},
		{
			name: "unable to read config file",
			args: args{
				file: "../../test/testdata/config/notexist.yaml",
			},
			want:       v1.Config{},
			wantErr:    true,
			wantErrMsg: "unable to read configuration file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("FindGoModuleAndPath() gotErrMsg = %v, wantErrMsg = %v", err.Error(), tt.wantErrMsg)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
