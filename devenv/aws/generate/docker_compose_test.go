package generate

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"syncgen/internal/config"
)

func TestWriteDockerComposeFile_Success(t *testing.T) {
	// Prepare a minimal template
	tmplText := `db:
  image: postgres
  environment:
	POSTGRES_DB: {{.DbName}}
	POSTGRES_USER: {{.DbUser}}
	POSTGRES_PASSWORD: {{.DbPassword}}
	WAL_LEVEL: {{.WalLevel}}
	MAX_WAL_SENDERS: {{.MaxWalSenders}}
	WAL_KEEP_SIZE: {{.WalKeepSize}}
	HOT_STANDBY: {{.HotStandby}}
	SYNCHRONOUS_COMMIT: {{.SynchronousCommit}}
`
	tmpl := template.Must(template.New("test").Parse(tmplText))

	cfg := &config.Config{
		Primary: config.Primary{
			DbName:     "testdb",
			DbUser:     "testuser",
			DbPassword: "testpass",
		},
		Options: config.Options{
			WalLevel:          "replica",
			MaxWalSenders:     5,
			WalKeepSize:       "2GB",
			HotStandby:        true,
			SynchronousCommit: "remote_apply",
		},
	}

	f, err := os.CreateTemp(t.TempDir(), "docker-compose.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	// Use the file for writing
	err = writeDockerComposeFile(tmpl, cfg, f)
	if err != nil {
		t.Fatalf("writeDockerComposeFile returned error: %v", err)
	}

	// Read back the file and check contents
	f.Seek(0, 0)
	content, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	// Check that all fields are rendered
	wantFields := []string{
		"POSTGRES_DB: testdb",
		"POSTGRES_USER: testuser",
		"POSTGRES_PASSWORD: testpass",
		"WAL_LEVEL: replica",
		"MAX_WAL_SENDERS: 5",
		"WAL_KEEP_SIZE: 2GB",
		"HOT_STANDBY: on",
		"SYNCHRONOUS_COMMIT: remote_apply",
	}
	for _, field := range wantFields {
		if !bytes.Contains(content, []byte(field)) {
			t.Errorf("output missing field: %s", field)
		}
	}
}

func TestWriteDockerComposeFile_TemplateError(t *testing.T) {
	// Template references a missing field
	tmplText := `{{.NonExistentField}}`
	tmpl := template.Must(template.New("bad").Parse(tmplText))

	cfg := &config.Config{
		Primary: config.Primary{
			DbName:     "testdb",
			DbUser:     "testuser",
			DbPassword: "testpass",
		},
		Options: config.Options{
			WalLevel:          "replica",
			MaxWalSenders:     5,
			WalKeepSize:       "2GB",
			HotStandby:        true,
			SynchronousCommit: "remote_apply",
		},
	}

	f, err := os.CreateTemp(t.TempDir(), "docker-compose.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	err = writeDockerComposeFile(tmpl, cfg, f)
	if err == nil {
		t.Error("expected error due to missing template field, got nil")
	}
}

func TestWriteDockerComposeFile_NilFile(t *testing.T) {
	tmpl := template.Must(template.New("empty").Parse(`test`))
	cfg := &config.Config{}
	err := writeDockerComposeFile(tmpl, cfg, nil)
	if err == nil {
		t.Error("expected error when file is nil, got nil")
	}
}

func TestWriteDockerComposeFile_EmptyTemplate(t *testing.T) {
	tmpl := template.Must(template.New("empty").Parse(""))
	cfg := &config.Config{}
	f, err := os.CreateTemp(t.TempDir(), "docker-compose.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()
	err = writeDockerComposeFile(tmpl, cfg, f)
	if err != nil {
		t.Errorf("expected no error for empty template, got: %v", err)
	}
}

func TestBoolToString(t *testing.T) {
	tests := []struct {
		input    bool
		expected string
	}{
		{true, "on"},
		{false, "off"},
	}

	for _, tt := range tests {
		got := boolToString(tt.input)
		if got != tt.expected {
			t.Errorf("boolToString(%v) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
