package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Primary struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Replica struct {
	Host         string `yaml:"host"`
	SyncInterval string `yaml:"sync_interval"`
}

type Options struct {
	RsyncUser        string `yaml:"rsync_user"`
	PromoteOnFailure bool   `yaml:"promote_on_failure"`
	Observability    string `yaml:"observability,omitempty"`
}

type Config struct {
	Primary  Primary   `yaml:"primary"`
	Replicas []Replica `yaml:"replicas"`
	Options  Options   `yaml:"options"`
}

func Parse(filename string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := Validate(&config); err != nil {
		return nil, fmt.Errorf("validation error:\n%w", err)
	}

	return &config, nil
}
