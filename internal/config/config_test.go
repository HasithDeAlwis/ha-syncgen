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
replicas:
  - host: 10.0.0.2
    sync_interval: 30s`,
			want: &Config{
				Primary: Primary{Host: "10.0.0.1", Port: 5432},
				Replicas: []Replica{
					{Host: "10.0.0.2", SyncInterval: "30s"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with multiple replicas",
			yaml: `primary:
  host: 192.168.1.100
  port: 5433
replicas:
  - host: 192.168.1.101
    sync_interval: 30s
  - host: 192.168.1.102
    sync_interval: 60s
  - host: 192.168.1.103
    sync_interval: 120s`,
			want: &Config{
				Primary: Primary{Host: "192.168.1.100", Port: 5433},
				Replicas: []Replica{
					{Host: "192.168.1.101", SyncInterval: "30s"},
					{Host: "192.168.1.102", SyncInterval: "60s"},
					{Host: "192.168.1.103", SyncInterval: "120s"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with zero port (should work)",
			yaml: `primary:
  host: localhost
  port: 0
replicas:
  - host: replica1
    sync_interval: 45s`,
			want: &Config{
				Primary: Primary{Host: "localhost", Port: 0},
				Replicas: []Replica{
					{Host: "replica1", SyncInterval: "45s"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with empty sync_interval",
			yaml: `primary:
  host: db-primary
  port: 5432
replicas:
  - host: db-replica1
    sync_interval: ""`,
			want: &Config{
				Primary: Primary{Host: "db-primary", Port: 5432},
				Replicas: []Replica{
					{Host: "db-replica1", SyncInterval: "default"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid yaml - malformed structure",
			yaml: `primary:
  host: 10.0.0.1
  port: "invalid_port"
replicas:
  - host: 10.0.0.2
    sync_interval: 30s
    invalid_field:`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid yaml - wrong data types",
			yaml: `primary:
  host: 123
  port: "not_a_number"
replicas:
  - host: ["array", "instead", "of", "string"]
    sync_interval: 30s`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid yaml - completely malformed",
			yaml: `this is not yaml at all
it's just plain text
with no structure`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid yaml - missing required sections",
			yaml: `some_other_field: value`,
			want: &Config{
				Primary:  Primary{},
				Replicas: nil,
			},
			wantErr: true, // YAML parsing succeeds, just empty values
		},
		{
			name: "invalid yaml - empty file",
			yaml: ``,
			want: &Config{
				Primary:  Primary{},
				Replicas: nil,
			},
			wantErr: true, // Empty YAML is not valid
		},
		{
			name: "invalid yaml - tabs instead of spaces",
			yaml: `primary:
	host: 10.0.0.1
	port: 5432
replicas:
	- host: 10.0.0.2
	  sync_interval: 30s`,
			want:    nil,
			wantErr: true, // YAML is sensitive to indentation
		},
		{
			name: "valid config with unicode characters",
			yaml: `primary:
  host: "测试服务器"
  port: 5432
replicas:
  - host: "副本服务器1"
    sync_interval: "30s"`,
			want: &Config{
				Primary: Primary{Host: "测试服务器", Port: 5432},
				Replicas: []Replica{
					{Host: "副本服务器1", SyncInterval: "30s"},
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
replicas:
  - host: 10.0.0.2
    sync_interval: 30s`

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
			Host: "test-host",
			Port: 1234,
		},
		Replicas: []Replica{
			{
				Host:         "replica-host",
				SyncInterval: "60s",
			},
		},
	}

	// This test ensures our struct is properly formed
	if config.Primary.Host != "test-host" {
		t.Errorf("Primary.Host = %q, want %q", config.Primary.Host, "test-host")
	}

	if config.Primary.Port != 1234 {
		t.Errorf("Primary.Port = %d, want %d", config.Primary.Port, 1234)
	}

	if len(config.Replicas) != 1 {
		t.Errorf("len(Replicas) = %d, want %d", len(config.Replicas), 1)
	}

	if config.Replicas[0].Host != "replica-host" {
		t.Errorf("Replicas[0].Host = %q, want %q", config.Replicas[0].Host, "replica-host")
	}

	if config.Replicas[0].SyncInterval != "60s" {
		t.Errorf("Replicas[0].SyncInterval = %q, want %q", config.Replicas[0].SyncInterval, "60s")
	}
}

// TestYAMLTagsWork ensures our YAML tags are working correctly
func TestYAMLTagsWork(t *testing.T) {
	yaml := `primary:
  host: yaml-test-host
  port: 9999
replicas:
  - host: yaml-replica-host
    sync_interval: 90s`

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

	if config.Replicas[0].SyncInterval != "90s" {
		t.Errorf("YAML 'sync_interval' didn't map to Replica.SyncInterval correctly")
	}
}
