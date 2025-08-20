package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper to create a mock template file in a temp directory
func setupMockTemplate(t *testing.T, content string) (string, string) {
	t.Helper()
	tmpDir := t.TempDir()
	tmplPath := filepath.Join(tmpDir, "mock.tmpl")
	if err := os.WriteFile(tmplPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write mock template: %v", err)
	}
	return tmpDir, tmplPath
}

// Helper to create a minimal DeploymentData struct
func mockDeploymentData() DeploymentData {
	return DeploymentData{
		PrimaryUser:     "tester",
		PrimaryIP:       "127.0.0.1",
		PrimaryPassword: "pw",
		PrimaryDBName:   "testdb",
		SSHKeyPath:      "/tmp/key.pem",
		SSHUser:         "ec2-user",
		Replicas: []ReplicaDeploymentData{
			{IP: "127.0.0.2", User: "replica", Password: "pw2", Name: "replica1"},
			{IP: "127.0.0.3", User: "replica2", Password: "pw3", Name: "replica2"},
		},
	}
}

func TestGenerateScriptFromTemplate(t *testing.T) {
	tmpl := "Hello, {{.PrimaryUser}}!"
	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	data := mockDeploymentData()
	err = generateScriptFromTemplate(tmpl, data, tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Hello, tester!" {
		t.Errorf("unexpected content: %s", content)
	}
}

func TestGenerateScriptFromTemplate_Replicas(t *testing.T) {
	tmpl := "Primary: {{.PrimaryUser}}@{{.PrimaryIP}}\nReplicas:\n{{range .Replicas}}- {{.Name}}: {{.User}}@{{.IP}} (pw:{{.Password}})\n{{end}}"
	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	data := mockDeploymentData()
	err = generateScriptFromTemplate(tmpl, data, tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	expected := "Primary: tester@127.0.0.1\nReplicas:\n- replica1: replica@127.0.0.2 (pw:pw2)\n- replica2: replica2@127.0.0.3 (pw:pw3)\n"
	if string(content) != expected {
		t.Errorf("unexpected content:\n%s\nwanted:\n%s", content, expected)
	}
}
