package config

import (
	"reflect"
	"strings"
	"testing"

	v1alpha1 "github.com/go-imports-organizer/goio/pkg/api/v1alpha1"
)

func TestLoad(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name       string
		args       args
		want       v1alpha1.Config
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "working config file",
			args: args{
				file: "../../test/testdata/config/works.yaml",
			},
			want: v1alpha1.Config{
				Excludes: []v1alpha1.Exclude{
					{
						MatchType: "name",
						RegExp:    "^\\.git$",
					},
					{
						MatchType: "name",
						RegExp:    "^vendor$",
					},
				},
				Groups: []v1alpha1.Group{
					{
						MatchOrder:  0,
						Description: "module",
						RegExp:      []string{"%{module}%"},
					},
					{
						MatchOrder:  1,
						Description: "standard",
						RegExp:      []string{"^[a-zA-Z0-9\\/]+$"},
					},
					{
						MatchOrder:  2,
						Description: "other",
						RegExp:      []string{"[a-zA-Z0-9]+\\.[a-zA-Z0-9]+/"},
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
			want:       v1alpha1.Config{},
			wantErr:    true,
			wantErrMsg: "unable to unmarshal file",
		},
		{
			name: "unable to read config file",
			args: args{
				file: "../../test/testdata/config/notexist.yaml",
			},
			want:       v1alpha1.Config{},
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
