package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	v1 "github.com/go-imports-organizer/gio/pkg/api/v1"
)

func Load(file string) (v1.Config, error) {
	var configFile []byte
	var err error
	if configFile, err = os.ReadFile(file); err != nil {
		return v1.Config{}, fmt.Errorf("unable to read configuration file %s: %s", file, err.Error())
	}

	var config v1.Config
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		return v1.Config{}, fmt.Errorf("unable to unmarshal file %s: %s", file, err.Error())
	}

	return config, nil
}
