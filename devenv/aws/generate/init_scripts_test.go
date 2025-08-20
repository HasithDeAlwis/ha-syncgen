package generate

import (
	"os"
	"path/filepath"
	"testing"

	"syncgen/internal/config"
)

func mockConfig() *config.Config {
	return &config.Config{
		Primary: config.Primary{
			DbUser:              "primaryuser",
			DbPassword:          "primarypw",
			ReplicationUser:     "repluser",
			ReplicationPassword: "replpw",
		},
		Replicas: []config.Replica{
			{
				DbUser:          "replicauser1",
				DbPassword:      "replicapw1",
				ReplicationSlot: "slot1",
			},
			{
				DbUser:          "replicauser2",
				DbPassword:      "replicapw2",
				ReplicationSlot: "slot2",
			},
		},
		Options: config.Options{
			WalLevel:          "logical",
			MaxWalSenders:     5,
			WalKeepSize:       "128MB",
			HotStandby:        true,
			SynchronousCommit: "on",
		},
		Monitoring: &config.Monitoring{
			Datadog: config.DatadogConfig{
				DatadogUserPassword: "ddpw",
			},
		},
	}
}

func TestLoadInitScriptTemplate(t *testing.T) {
	tempRoot := t.TempDir()
	generateDir := filepath.Join(tempRoot, "generate")
	templatesDir := filepath.Join(generateDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}

	tmplName := "test-init.sql.tmpl"
	tmplContent := "-- SQL for {{.DbUser}}"
	tmplPath := filepath.Join(templatesDir, tmplName)
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	if err := os.Chdir(generateDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origWD)

	tmpl, err := loadInitScriptTemplate(tmplName)
	if err != nil {
		t.Fatalf("failed to load template: %v", err)
	}
	if tmpl == nil {
		t.Fatal("template is nil")
	}
}

func TestGeneratePrimaryInitScript(t *testing.T) {
	tempRoot := t.TempDir()
	generateDir := filepath.Join(tempRoot, "generate")
	templatesDir := filepath.Join(generateDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}

	tmplContent := "-- PRIMARY: {{.DbUser}} {{.DbPassword}} {{.ReplicationUser}} {{.ReplicationPassword}} {{.DatadogPassword}} {{.WalLevel}} {{.MaxWalSenders}} {{.WalKeepSize}} {{.HotStandby}} {{.SynchronousCommit}}"
	tmplPath := filepath.Join(templatesDir, "primary-init.sql.tmpl")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	if err := os.Chdir(generateDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origWD)

	cfg := mockConfig()
	err = GeneratePrimaryInitScript(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outputFile := filepath.Join(tempRoot, "generated", "primary", "init-scripts", "01-setup-primary.sql")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	expected := "-- PRIMARY: primaryuser primarypw repluser replpw ddpw logical 5 128MB true on"
	if string(content) != expected {
		t.Errorf("unexpected content: %s", content)
	}
}

func TestGenerateReplicaInitScripts(t *testing.T) {
	tempRoot := t.TempDir()
	generateDir := filepath.Join(tempRoot, "generate")
	templatesDir := filepath.Join(generateDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}

	tmplContent := "-- REPLICA: {{.DbUser}} {{.DbPassword}} {{.DatadogPassword}} {{.HotStandby}} {{.ReplicationSlot}}"
	tmplPath := filepath.Join(templatesDir, "replica-init.sql.tmpl")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	if err := os.Chdir(generateDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origWD)

	cfg := mockConfig()
	err = GenerateReplicaInitScripts(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, replica := range cfg.Replicas {
		replicaName := "replica" + string(rune('1'+i))
		outputFile := filepath.Join(tempRoot, "generated", replicaName, "init-scripts", "01-setup-"+replicaName+".sql")
		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("failed to read output for %s: %v", replicaName, err)
		}
		expected := "-- REPLICA: " + replica.DbUser + " " + replica.DbPassword + " ddpw true " + replica.ReplicationSlot
		if string(content) != expected {
			t.Errorf("unexpected content for %s: %s", replicaName, content)
		}
	}
}
