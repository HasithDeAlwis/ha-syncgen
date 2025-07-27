package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config with single replica",
			yaml: `primary:
  host: 10.0.0.1
  port: 5432
  data_directory: /var/lib/postgresql/data
  replication_user: replicator
  replication_password: secure_password
replicas:
  - host: 10.0.0.2
    port: 5432
    replication_slot: replica_slot_1
    sync_mode: async
options:
  promote_on_failure: true
  wal_level: replica
  max_wal_senders: 3
  wal_keep_size: 1GB
  hot_standby: true
  synchronous_commit: on`,
			want: &Config{
				Primary: Primary{
					Host:                "10.0.0.1",
					Port:                5432,
					DataDirectory:       "/var/lib/postgresql/data",
					ReplicationUser:     "replicator",
					ReplicationPassword: "secure_password",
				},
				Replicas: []Replica{
					{
						Host:            "10.0.0.2",
						Port:            5432,
						ReplicationSlot: "replica_slot_1",
						SyncMode:        "async",
					},
				},
				Options: Options{
					PromoteOnFailure:  true,
					WalLevel:          "replica",
					MaxWalSenders:     3,
					WalKeepSize:       "1GB",
					HotStandby:        true,
					SynchronousCommit: "on",
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with multiple replicas",
			yaml: `primary:
  host: 192.168.1.100
  port: 5433
  data_directory: /opt/postgresql/data
  replication_user: repl_user
  replication_password: secret123
replicas:
  - host: 192.168.1.101
    port: 5433
    replication_slot: replica_1
    sync_mode: sync
  - host: 192.168.1.102
    port: 5433
    replication_slot: replica_2
    sync_mode: async
  - host: 192.168.1.103
    port: 5433
    replication_slot: replica_3
    sync_mode: async
options:
  promote_on_failure: false
  wal_level: logical
  max_wal_senders: 5
  wal_keep_size: 2GB
  hot_standby: true
  synchronous_commit: remote_apply`,
			want: &Config{
				Primary: Primary{
					Host:                "192.168.1.100",
					Port:                5433,
					DataDirectory:       "/opt/postgresql/data",
					ReplicationUser:     "repl_user",
					ReplicationPassword: "secret123",
				},
				Replicas: []Replica{
					{Host: "192.168.1.101", Port: 5433, ReplicationSlot: "replica_1", SyncMode: "sync"},
					{Host: "192.168.1.102", Port: 5433, ReplicationSlot: "replica_2", SyncMode: "async"},
					{Host: "192.168.1.103", Port: 5433, ReplicationSlot: "replica_3", SyncMode: "async"},
				},
				Options: Options{
					PromoteOnFailure:  false,
					WalLevel:          "logical",
					MaxWalSenders:     5,
					WalKeepSize:       "2GB",
					HotStandby:        true,
					SynchronousCommit: "remote_apply",
				},
			},
			wantErr: false,
		},
		{
			name: "config with defaults applied",
			yaml: `primary:
  host: localhost
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1`,
			want: &Config{
				Primary: Primary{
					Host:                "localhost",
					Port:                5432,
					DataDirectory:       "/var/lib/postgresql/data",
					ReplicationUser:     "replicator",
					ReplicationPassword: "password",
				},
				Replicas: []Replica{
					{
						Host:            "replica1",
						Port:            5432,
						ReplicationSlot: "slot1",
						SyncMode:        "async",
					},
				},
				Options: Options{
					PromoteOnFailure:  false,
					WalLevel:          "replica",
					MaxWalSenders:     3,
					WalKeepSize:       "1GB",
					HotStandby:        false,
					SynchronousCommit: "on",
				},
			},
			wantErr: false,
		},
		{
			name: "missing primary host",
			yaml: `primary:
  port: 5432
  data_directory: /var/lib/postgresql/data
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing replication user",
			yaml: `primary:
  host: primary
  port: 5432
  data_directory: /var/lib/postgresql/data
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing replication password",
			yaml: `primary:
  host: primary
  port: 5432
  data_directory: /var/lib/postgresql/data
  replication_user: replicator
replicas:
  - host: replica1
    replication_slot: slot1`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "duplicate replication slots",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: duplicate_slot
  - host: replica2
    replication_slot: duplicate_slot`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid sync mode",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1
    sync_mode: invalid_mode`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid wal level",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1
options:
  wal_level: invalid_level`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid replication slot name",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: "invalid-slot-name-with-dashes"`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid wal keep size",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1
options:
  wal_keep_size: "invalid_size"`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "no replicas",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas: []`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid yaml - malformed structure",
			yaml: `primary:
  host: 10.0.0.1
  port: "invalid_port"
replicas:
  - host: 10.0.0.2
    replication_slot: slot1
    invalid_field:`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid config with unicode characters",
			yaml: `primary:
  host: "测试服务器"
  port: 5432
  data_directory: /var/lib/postgresql/data
  replication_user: replicator
  replication_password: password
replicas:
  - host: "副本服务器1"
    port: 5432
    replication_slot: "replica_slot_1"
    sync_mode: async`,
			want: &Config{
				Primary: Primary{
					Host:                "测试服务器",
					Port:                5432,
					DataDirectory:       "/var/lib/postgresql/data",
					ReplicationUser:     "replicator",
					ReplicationPassword: "password",
				},
				Replicas: []Replica{
					{
						Host:            "副本服务器1",
						Port:            5432,
						ReplicationSlot: "replica_slot_1",
						SyncMode:        "async",
					},
				},
				Options: Options{
					PromoteOnFailure:  false,
					WalLevel:          "replica",
					MaxWalSenders:     3,
					WalKeepSize:       "1GB",
					HotStandby:        false,
					SynchronousCommit: "on",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(tmpFile, []byte(tt.yaml), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test the Parse function
			got, err := Parse(tmpFile)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, don't check the result
			if tt.wantErr {
				return
			}

			// Check result
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseFileNotFound(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "nonexistent file",
			filename: "nonexistent.yaml",
		},
		{
			name:     "file in nonexistent directory",
			filename: "/nonexistent/directory/config.yaml",
		},
		{
			name:     "empty filename",
			filename: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.filename)
			if err == nil {
				t.Errorf("Parse() expected error for filename %q, but got none", tt.filename)
			}

			// Check that error message contains helpful information
			if !strings.Contains(err.Error(), "failed to read config file") {
				t.Errorf("Parse() error should mention file reading failure, got: %v", err)
			}
		})
	}
}

func TestParseFilePermissions(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	validYAML := `primary:
  host: 10.0.0.1
  port: 5432
  data_directory: /var/lib/postgresql/data
  replication_user: replicator
  replication_password: password
replicas:
  - host: 10.0.0.2
    port: 5432
    replication_slot: slot1
    sync_mode: async`

	err := os.WriteFile(tmpFile, []byte(validYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Remove read permissions
	err = os.Chmod(tmpFile, 0000)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}

	// Restore permissions after test
	defer func() {
		os.Chmod(tmpFile, 0644)
	}()

	// Test parsing the file without read permissions
	_, err = Parse(tmpFile)
	if err == nil {
		t.Error("Parse() expected error for file without read permissions, but got none")
	}
}

func TestParseErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		yaml          string
		expectInError string
	}{
		{
			name: "yaml unmarshal error",
			yaml: `primary:
  host: 10.0.0.1
  port: [invalid, array, for, port]`,
			expectInError: "cannot unmarshal",
		},
		{
			name: "malformed yaml structure",
			yaml: `primary
  host: 10.0.0.1
  port: 5432`,
			expectInError: "mapping values are not allowed in this context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(tmpFile, []byte(tt.yaml), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			_, err = Parse(tmpFile)
			if err == nil {
				t.Error("Parse() expected error but got none")
				return
			}

			if !strings.Contains(err.Error(), tt.expectInError) {
				t.Errorf("Parse() error should contain %q, got: %v", tt.expectInError, err)
			}
		})
	}
}

func TestConfigStructFields(t *testing.T) {
	// Test that struct fields are properly tagged for YAML
	config := &Config{
		Primary: Primary{
			Host:                "test-host",
			Port:                1234,
			DataDirectory:       "/test/data",
			ReplicationUser:     "test_user",
			ReplicationPassword: "test_password",
		},
		Replicas: []Replica{
			{
				Host:            "replica-host",
				Port:            5432,
				ReplicationSlot: "test_slot",
				SyncMode:        "async",
			},
		},
		Options: Options{
			PromoteOnFailure:  true,
			WalLevel:          "replica",
			MaxWalSenders:     3,
			WalKeepSize:       "1GB",
			HotStandby:        true,
			SynchronousCommit: "on",
		},
	}

	// Test Primary fields
	if config.Primary.Host != "test-host" {
		t.Errorf("Primary.Host = %q, want %q", config.Primary.Host, "test-host")
	}
	if config.Primary.Port != 1234 {
		t.Errorf("Primary.Port = %d, want %d", config.Primary.Port, 1234)
	}
	if config.Primary.ReplicationUser != "test_user" {
		t.Errorf("Primary.ReplicationUser = %q, want %q", config.Primary.ReplicationUser, "test_user")
	}

	// Test Replica fields
	if len(config.Replicas) != 1 {
		t.Errorf("len(Replicas) = %d, want %d", len(config.Replicas), 1)
	}
	if config.Replicas[0].Host != "replica-host" {
		t.Errorf("Replicas[0].Host = %q, want %q", config.Replicas[0].Host, "replica-host")
	}
	if config.Replicas[0].ReplicationSlot != "test_slot" {
		t.Errorf("Replicas[0].ReplicationSlot = %q, want %q", config.Replicas[0].ReplicationSlot, "test_slot")
	}
	if config.Replicas[0].SyncMode != "async" {
		t.Errorf("Replicas[0].SyncMode = %q, want %q", config.Replicas[0].SyncMode, "async")
	}

	// Test Options fields
	if config.Options.WalLevel != "replica" {
		t.Errorf("Options.WalLevel = %q, want %q", config.Options.WalLevel, "replica")
	}
	if config.Options.MaxWalSenders != 3 {
		t.Errorf("Options.MaxWalSenders = %d, want %d", config.Options.MaxWalSenders, 3)
	}
}

// TestYAMLTagsWork ensures our YAML tags are working correctly
func TestYAMLTagsWork(t *testing.T) {
	yaml := `primary:
  host: yaml-test-host
  port: 9999
  data_directory: /yaml/test/data
  replication_user: yaml_user
  replication_password: yaml_password
replicas:
  - host: yaml-replica-host
    port: 5433
    replication_slot: yaml_slot
    sync_mode: sync
options:
  promote_on_failure: true
  wal_level: logical
  max_wal_senders: 5
  wal_keep_size: 2GB
  hot_standby: true
  synchronous_commit: remote_apply`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(tmpFile, []byte(yaml), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config, err := Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Test that YAML fields map correctly to struct fields
	if config.Primary.Host != "yaml-test-host" {
		t.Errorf("YAML 'host' didn't map to Primary.Host correctly")
	}
	if config.Primary.Port != 9999 {
		t.Errorf("YAML 'port' didn't map to Primary.Port correctly")
	}
	if config.Primary.ReplicationUser != "yaml_user" {
		t.Errorf("YAML 'replication_user' didn't map to Primary.ReplicationUser correctly")
	}
	if config.Replicas[0].ReplicationSlot != "yaml_slot" {
		t.Errorf("YAML 'replication_slot' didn't map to Replica.ReplicationSlot correctly")
	}
	if config.Replicas[0].SyncMode != "sync" {
		t.Errorf("YAML 'sync_mode' didn't map to Replica.SyncMode correctly")
	}
	if config.Options.WalLevel != "logical" {
		t.Errorf("YAML 'wal_level' didn't map to Options.WalLevel correctly")
	}
	if config.Options.MaxWalSenders != 5 {
		t.Errorf("YAML 'max_wal_senders' didn't map to Options.MaxWalSenders correctly")
	}
}

// TestValidationCases tests specific validation scenarios
func TestValidationCases(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid minimal config",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1`,
			wantErr: false,
		},
		{
			name: "duplicate replication slots",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: same_slot
  - host: replica2
    replication_slot: same_slot`,
			wantErr: true,
			errMsg:  "already used",
		},
		{
			name: "invalid replication slot name with dashes",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: invalid-slot-name`,
			wantErr: true,
			errMsg:  "invalid replication slot name",
		},
		{
			name: "invalid sync mode",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1
    sync_mode: invalid`,
			wantErr: true,
			errMsg:  "invalid sync_mode",
		},
		{
			name: "invalid wal level",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1
options:
  wal_level: invalid`,
			wantErr: true,
			errMsg:  "invalid wal_level",
		},
		{
			name: "invalid wal keep size",
			yaml: `primary:
  host: primary
  replication_user: replicator
  replication_password: password
replicas:
  - host: replica1
    replication_slot: slot1
options:
  wal_keep_size: invalid_size_format`,
			wantErr: true,
			errMsg:  "invalid wal_keep_size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(tmpFile, []byte(tt.yaml), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			_, err = Parse(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Parse() error should contain %q, got: %v", tt.errMsg, err)
			}
		})
	}
}
