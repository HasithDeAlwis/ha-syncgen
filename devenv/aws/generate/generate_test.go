package generate

import (
	"reflect"
	"testing"

	"syncgen/internal/config"
)

func TestParseJSON_Valid(t *testing.T) {
	tests := []struct {
		jsonConfig   string
		wantedConfig *config.Config
	}{
		{
			jsonConfig: `{
				"outputs": {
					"datadog_details": {
						"value": {
							"api_key": "your_datadog_api_key",
							"site": "datadoghq.com"
						}
					},
					"instance_details": {
						"value": {
							"primary": {
								"db_name": "test_db",
								"db_user": "test_user",
								"db_password": "test_password",
								"ip_address": "192.168.1.1",
								"role": "primary"
							},
							"replica": {
								"db_user": "replica_user",
								"db_password": "replica_password",
								"ip_address": "192.168.1.2",
								"role": "replica"
							}
						}
					}
				}
			}`,
			wantedConfig: &config.Config{
				Primary: config.Primary{
					Host:                "192.168.1.1",
					DbName:              "test_db",
					DbUser:              "test_user",
					DbPassword:          "test_password",
					ReplicationUser:     "primary_replica_user",
					Port:                5432,
					ReplicationPassword: "replica_user_password",
					DataDirectory:       "/var/lib/postgresql/data/primary",
				},
				Replicas: []config.Replica{
					{
						Host:            "192.168.1.2",
						DbUser:          "replica_user",
						DbPassword:      "replica_password",
						Port:            5432,
						SyncMode:        "async",
						ReplicationSlot: "replica_slot_1",
					},
				},
				Options: config.Options{
					PromoteOnFailure:  true,
					WalLevel:          "replica",
					MaxWalSenders:     3,
					WalKeepSize:       "1GB",
					HotStandby:        true,
					SynchronousCommit: "on",
				},
				Monitoring: &config.Monitoring{
					Datadog: config.DatadogConfig{
						Enabled:             true,
						ApiKey:              "your_datadog_api_key",
						Site:                "datadoghq.com",
						DatadogUserPassword: "test_password",
					},
				},
			},
		},
		{
			jsonConfig: `{
				"outputs": {
					"datadog_details": {
						"value": {
							"api_key": "test_api_key",
							"site": "datadoghq.com"
						}
					},
					"instance_details": {
						"value": {
							"primary": {
								"db_name": "test_db",
								"db_user": "test_user",
								"db_password": "test_password",
								"ip_address": "192.168.1.1",
								"role": "primary"
							},
							"replica1": {
								"db_user": "replica1_user",
								"db_password": "replica1_password",
								"ip_address": "192.168.1.2",
								"role": "replica"
							},
							"replica2": {
								"db_user": "replica2_user",
								"db_password": "replica2_password",
								"ip_address": "192.168.1.3",
								"role": "replica"
							}
						}
					}
				}
			}`,
			wantedConfig: &config.Config{
				Primary: config.Primary{
					DbName:              "test_db",
					DbUser:              "test_user",
					DbPassword:          "test_password",
					Host:                "192.168.1.1",
					Port:                5432,
					ReplicationUser:     "primary_replica_user",
					ReplicationPassword: "replica_user_password",
					DataDirectory:       "/var/lib/postgresql/data/primary",
				},
				Replicas: []config.Replica{
					{
						DbUser:          "replica1_user",
						DbPassword:      "replica1_password",
						Host:            "192.168.1.2",
						Port:            5432,
						ReplicationSlot: "replica_slot_1",
						SyncMode:        "async",
					},
					{
						DbUser:          "replica2_user",
						DbPassword:      "replica2_password",
						Host:            "192.168.1.3",
						Port:            5432,
						ReplicationSlot: "replica_slot_2",
						SyncMode:        "async",
					},
				},
				Options: config.Options{
					PromoteOnFailure:  true,
					WalLevel:          "replica",
					MaxWalSenders:     3,
					WalKeepSize:       "1GB",
					HotStandby:        true,
					SynchronousCommit: "on",
				},
				Monitoring: &config.Monitoring{
					Datadog: config.DatadogConfig{
						Enabled:             true,
						ApiKey:              "test_api_key",
						Site:                "datadoghq.com",
						DatadogUserPassword: "test_password",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("Parse JSON", func(t *testing.T) {
			gotConfig, err := parseJSON([]byte(tt.jsonConfig))
			if err != nil {
				t.Errorf("ParseJSON() error = %v", err)
				return
			}

			if !reflect.DeepEqual(gotConfig, tt.wantedConfig) {
				t.Errorf("ParseJSON() = %v, want %v", gotConfig, tt.wantedConfig)
			}
		})
	}
}

func TestParseJSON_MalformedJSON(t *testing.T) {
	malformed := []byte(`{"outputs": { "datadog_details": { "value": { "api_key": "key", "site": "datadoghq.com" } }`) // missing closing braces
	_, err := parseJSON(malformed)
	if err == nil {
		t.Error("Expected error for malformed JSON, got nil")
	}
}

func TestParseJSON_EmptyInput(t *testing.T) {
	empty := []byte(``)
	_, err := parseJSON(empty)
	if err == nil {
		t.Error("Expected error for empty input, got nil")
	}
}

func TestParseJSON_MissingRequiredFields(t *testing.T) {
	missingFields := []byte(`{"outputs": {}}`)
	_, err := parseJSON(missingFields)
	if err == nil {
		t.Error("Expected error for missing required fields, got nil")
	}
}

func TestParseJSON_InvalidFieldValues(t *testing.T) {
	invalidValue := []byte(`{
		"outputs": {
			"datadog_details": {
				"value": {
					"api_key": "",
					"site": 123
				}
			}
		}
	}`)
	_, err := parseJSON(invalidValue)
	if err == nil {
		t.Error("Expected error for invalid field values, got nil")
	}
}
