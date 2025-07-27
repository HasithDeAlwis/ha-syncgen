package config

import (
	"errors"
	"fmt"
	"time"
)

func Validate(cfg *Config) error {
	var errs []error

	// Validate primary
	if cfg.Primary.Host == "" {
		errs = append(errs, fmt.Errorf("primary.host is required"))
	}

	// Validate replicas
	if len(cfg.Replicas) == 0 {
		errs = append(errs, fmt.Errorf("at least one replica is required"))
	}

	for i := range cfg.Replicas {
		if cfg.Replicas[i].Host == "" {
			errs = append(errs, fmt.Errorf("replicas[%d].host is required", i))
		}
		if cfg.Replicas[i].SyncInterval == "" {
			cfg.Replicas[i].SyncInterval = "default"
		} else if err := validateTimeInterval(cfg.Replicas[i].SyncInterval); err != nil {
			errs = append(errs, fmt.Errorf("replicas[%d].sync_interval: %w", i, err))
		}
	}

	// Validate options
	if cfg.Options.RsyncUser == "" {
		errs = append(errs, fmt.Errorf("options.rsync_user is required"))
	}

	return errors.Join(errs...)
}

func validateTimeInterval(interval string) error {
	if _, err := time.ParseDuration(interval); err != nil {
		return fmt.Errorf("invalid time interval: %w", err)
	}
	return nil
}
