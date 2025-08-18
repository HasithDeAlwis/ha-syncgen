package generate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syncgen/internal/config"

	yaml "gopkg.in/yaml.v3"
)

type DatadogValue struct {
	APIKey string `json:"api_key"`
	Site   string `json:"site"`
}

type DatadogItem struct {
	Value *DatadogValue `json:"value,omitempty"`
}

type Instance struct {
	DbName     string `json:"db_name,omitempty"`
	DbUser     string `json:"db_user,omitempty"`
	DbPassword string `json:"db_password,omitempty"`
	IPAddress  string `json:"ip_address,omitempty"`
	Role       string `json:"role,omitempty"`
}

type InstanceDetailsItem struct {
	Value map[string]Instance `json:"value,omitempty"`
}

type Outputs struct {
	DatadogDetails  *DatadogItem         `json:"datadog_details,omitempty"`
	InstanceDetails *InstanceDetailsItem `json:"instance_details,omitempty"`
}

type TFRoot struct {
	Outputs *Outputs `json:"outputs,omitempty"`
}

// PrintTFRoot prints the TFRoot struct in a pretty JSON format.
func PrintTFRoot(root *TFRoot) {
	if root == nil {
		fmt.Println("TFRoot is nil")
		return
	}
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling TFRoot: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func fromJSON(root *TFRoot) (*config.Config, error) {
	// Validate root and outputs presence
	if root == nil || root.Outputs == nil {
		return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
	}
	outs := root.Outputs

	cfg := &config.Config{}

	// Set sensible defaults for options
	cfg.Options = config.Options{
		PromoteOnFailure:  true,
		WalLevel:          "replica",
		MaxWalSenders:     3,
		WalKeepSize:       "1GB",
		HotStandby:        true,
		SynchronousCommit: "on",
	}

	// Require instance details
	if outs.InstanceDetails == nil || len(outs.InstanceDetails.Value) == 0 {
		return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
	}

	// Parse instances
	keys := make([]string, 0, len(outs.InstanceDetails.Value))
	for k := range outs.InstanceDetails.Value {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Find primary
	foundPrimary := false
	for _, name := range keys {
		inst := outs.InstanceDetails.Value[name]
		if strings.EqualFold(name, "primary") || strings.EqualFold(inst.Role, "primary") {
			// Validate required primary fields
			if inst.IPAddress == "" || inst.DbUser == "" || inst.DbPassword == "" {
				return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
			}
			cfg.Primary = config.Primary{
				Host:                inst.IPAddress,
				Port:                5432,
				DbName:              inst.DbName,
				DbUser:              inst.DbUser,
				DbPassword:          inst.DbPassword,
				ReplicationUser:     "primary_replica_user",
				ReplicationPassword: "replica_user_password",
				DataDirectory:       "/var/lib/postgresql/data/primary",
			}
			foundPrimary = true
			break
		}
	}
	if !foundPrimary {
		return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
	}

	// Collect replicas
	for i, name := range keys {
		inst := outs.InstanceDetails.Value[name]

		if strings.EqualFold(inst.Role, "replica") {
			// Validate required replica fields
			if inst.IPAddress == "" || inst.DbUser == "" || inst.DbPassword == "" {
				return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
			}
			rep := config.Replica{
				Host:            inst.IPAddress,
				Port:            5432,
				DbUser:          inst.DbUser,
				DbPassword:      inst.DbPassword,
				SyncMode:        "async",
				ReplicationSlot: "replica_slot_" + strconv.Itoa(i),
			}
			cfg.Replicas = append(cfg.Replicas, rep)
		}
	}

	// Require at least one replica
	if len(cfg.Replicas) == 0 {
		return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
	}

	var dd *DatadogValue
	if outs.DatadogDetails == nil || outs.DatadogDetails.Value == nil {
		return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
	}

	dd = outs.DatadogDetails.Value
	if dd == nil || dd.APIKey == "" || dd.Site == "" {
		return nil, fmt.Errorf("invalid JSON configuration generated during Terraform creation")
	}

	cfg.Monitoring = &config.Monitoring{Datadog: config.DatadogConfig{
		Enabled:             true,
		ApiKey:              dd.APIKey,
		Site:                dd.Site,
		DatadogUserPassword: cfg.Primary.DbPassword,
	}}

	return cfg, nil
}

func parseJSON(input []byte) (*config.Config, error) {
	var root TFRoot
	if err := json.Unmarshal(input, &root); err != nil {
		return nil, err
	}

	return fromJSON(&root)
}

func writeRelativeToSource(rel string, data []byte, perm os.FileMode) (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}

	absPath, err := filepath.Abs(filename)
	target := filepath.Join(filepath.Dir(absPath), rel)

	if err != nil {
		return "", err
	}

	if err := os.WriteFile(target, data, perm); err != nil {
		return "", err
	}

	return target, nil
}

func generateYAML(cfg *config.Config) ([]byte, error) {
	if cfg == nil {
		return []byte(""), fmt.Errorf("nil config provided")
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return []byte(""), err
	}

	return out, nil
}

func ParseTFOutputsFile(filepath string) (string, error) {
	fileData, readFileErr := readFile(filepath)
	if readFileErr != nil {
		return "", fmt.Errorf("failed to read file: %v", readFileErr)
	}

	cfg, parseJSONErr := parseJSON(fileData)
	if parseJSONErr != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", parseJSONErr)
	}

	yaml, generateYAMLErr := generateYAML(cfg)
	if generateYAMLErr != nil {
		return "", fmt.Errorf("failed to generate YAML: %v", generateYAMLErr)
	}

	absPath, writeToFileErr := writeRelativeToSource("../config.yaml", yaml, 0644)
	if writeToFileErr != nil {
		return "", fmt.Errorf("failed to write YAML file: %v", writeToFileErr)
	}

	return absPath, nil
}
