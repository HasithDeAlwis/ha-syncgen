package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Primary struct {
	Host                string `yaml:"host"`
	Port                int    `yaml:"port"`
	DataDirectory       string `yaml:"data_directory"`
	DbName              string `yaml:"db_name"`
	DbUser              string `yaml:"db_user"`
	DbPassword          string `yaml:"db_password"`
	ReplicationUser     string `yaml:"replication_user"`
	ReplicationPassword string `yaml:"replication_password"`
}

type Replica struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	ReplicationSlot string `yaml:"replication_slot"`
	SyncMode        string `yaml:"sync_mode"`
	DbUser          string `yaml:"db_user"`
	DbPassword      string `yaml:"db_password"`
}

type Options struct {
	PromoteOnFailure  bool   `yaml:"promote_on_failure"`
	WalLevel          string `yaml:"wal_level"`
	MaxWalSenders     int    `yaml:"max_wal_senders"`
	WalKeepSize       string `yaml:"wal_keep_size"`
	HotStandby        bool   `yaml:"hot_standby"`
	SynchronousCommit string `yaml:"synchronous_commit"`
}

type Monitoring struct {
	Datadog DatadogConfig `yaml:"datadog"`
}

type DatadogConfig struct {
	Enabled             bool   `yaml:"enabled"`
	ApiKey              string `yaml:"api_key"`
	Site                string `yaml:"site"`
	DatadogUserPassword string `yaml:"datadog_user_password"`
}

type Config struct {
	Primary    Primary     `yaml:"primary"`
	Replicas   []Replica   `yaml:"replicas"`
	Options    Options     `yaml:"options"`
	Monitoring *Monitoring `yaml:"monitoring,omitempty"`
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
