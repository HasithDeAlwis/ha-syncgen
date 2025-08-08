package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Validate(cfg *Config) error {
	var errs []error

	errs = append(errs, validatePrimaryConfig(cfg.Primary)...)

	if len(cfg.Replicas) == 0 {
		errs = append(errs, fmt.Errorf("at least one replica is required"))
	}
	replicationSlots := make(map[string]bool)
	for i := range cfg.Replicas {
		replica := &cfg.Replicas[i]
		errs = append(errs, validateReplicaConfig(*replica, replicationSlots, i)...)
	}

	errs = append(errs, validateOptions(&cfg.Options)...)

	return errors.Join(errs...)
}

func validateOptions(options *Options) []error {
	var errs []error

	if options.WalLevel == "" {
		options.WalLevel = "replica"
	} else if err := validateWalLevel(options.WalLevel); err != nil {
		errs = append(errs, fmt.Errorf("options.wal_level: %w", err))
	}

	if options.MaxWalSenders <= 0 {
		options.MaxWalSenders = 3 // Default
	}

	if options.WalKeepSize == "" {
		options.WalKeepSize = "1GB"
	} else if err := validateWalKeepSize(options.WalKeepSize); err != nil {
		errs = append(errs, fmt.Errorf("options.wal_keep_size: %w", err))
	}

	if options.SynchronousCommit == "" {
		options.SynchronousCommit = "on"
	} else if err := validateSynchronousCommit(options.SynchronousCommit); err != nil {
		errs = append(errs, fmt.Errorf("options.synchronous_commit: %w", err))
	}
	return errs
}

func validateReplicaConfig(replica Replica, replicationSlots map[string]bool, i int) []error {
	var errs []error

	if replica.Host == "" {
		errs = append(errs, fmt.Errorf("replicas[%d].host is required", i))
	}
	if replica.Port <= 0 {
		replica.Port = 5432 // Default PostgreSQL port
	}
	if replica.ReplicationSlot == "" {
		errs = append(errs, fmt.Errorf("replicas[%d].replication_slot is required", i))
	} else {
		// Check for duplicate replication slots
		if replicationSlots[replica.ReplicationSlot] {
			errs = append(errs, fmt.Errorf("replicas[%d].replication_slot '%s' is already used", i, replica.ReplicationSlot))
		}
		replicationSlots[replica.ReplicationSlot] = true

		// Validate replication slot name format
		if err := validateReplicationSlotName(replica.ReplicationSlot); err != nil {
			errs = append(errs, fmt.Errorf("replicas[%d].replication_slot: %w", i, err))
		}
	}

	if replica.SyncMode == "" {
		replica.SyncMode = "async" // Default to async
	} else if err := validateSyncMode(replica.SyncMode); err != nil {
		errs = append(errs, fmt.Errorf("replicas[%d].sync_mode: %w", i, err))
	}
	return errs
}

func validatePrimaryConfig(primary Primary) []error {
	var errs []error

	if primary.Host == "" {
		errs = append(errs, fmt.Errorf("primary.host is required"))
	}
	if primary.Port <= 0 {
		primary.Port = 5432 // Default PostgreSQL port
	}
	if primary.DataDirectory == "" {
		primary.DataDirectory = "/var/lib/postgresql/data" // Default
	}
	if primary.ReplicationUser == "" {
		errs = append(errs, fmt.Errorf("primary.replication_user is required"))
	}
	if primary.ReplicationPassword == "" {
		errs = append(errs, fmt.Errorf("primary.replication_password is required"))
	}

	if primary.DbName == "" {
		errs = append(errs, fmt.Errorf("primary.db_name is required"))
	}
	if primary.DbUser == "" {
		errs = append(errs, fmt.Errorf("primary.db_user is required"))
	}
	if primary.DbPassword == "" {
		errs = append(errs, fmt.Errorf("primary.db_password is required"))
	}

	return errs
}

func validateReplicationSlotName(name string) error {
	// PostgreSQL replication slot names must be valid SQL identifiers
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	if !matched {
		return fmt.Errorf("invalid replication slot name '%s': must be a valid SQL identifier", name)
	}
	if len(name) > 63 {
		return fmt.Errorf("replication slot name '%s' too long: maximum 63 characters", name)
	}
	return nil
}

func validateSyncMode(mode string) error {
	validModes := []string{"sync", "async"}
	for _, valid := range validModes {
		if mode == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid sync_mode '%s': must be one of %v", mode, validModes)
}

func validateWalLevel(level string) error {
	validLevels := []string{"minimal", "replica", "logical"}
	for _, valid := range validLevels {
		if level == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid wal_level '%s': must be one of %v", level, validLevels)
}

func validateWalKeepSize(size string) error {
	// Accept formats like "1GB", "512MB", "2048", etc.
	if size == "0" {
		return nil
	}

	// Check if it's just a number (in MB)
	if _, err := strconv.Atoi(size); err == nil {
		return nil
	}

	// Check if it has valid unit suffix
	validUnits := []string{"kB", "MB", "GB", "TB"}
	for _, unit := range validUnits {
		if strings.HasSuffix(size, unit) {
			numPart := strings.TrimSuffix(size, unit)
			if _, err := strconv.Atoi(numPart); err != nil {
				return fmt.Errorf("invalid wal_keep_size '%s': numeric part must be an integer", size)
			}
			return nil
		}
	}

	return fmt.Errorf("invalid wal_keep_size '%s': must be a number with optional unit (kB, MB, GB, TB)", size)
}

func validateSynchronousCommit(commit string) error {
	validValues := []string{"on", "off", "local", "remote_write", "remote_apply"}
	for _, valid := range validValues {
		if commit == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid synchronous_commit '%s': must be one of %v", commit, validValues)
}
