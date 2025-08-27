package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"syncgen/internal/config"
)

// InitScriptData holds the data for SQL init script template rendering
type InitScriptData struct {
	DbUser              string
	DbPassword          string
	ReplicationUser     string
	ReplicationPassword string
	DatadogPassword     string
	WalLevel            string
	MaxWalSenders       int
	WalKeepSize         string
	HotStandby          string
	SynchronousCommit   string
	ReplicationSlot     string
}

func loadInitScriptTemplate(templateName string) (*template.Template, error) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("unable to get caller information")
	}

	absPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return nil, err
	}

	generateDir := filepath.Dir(absPath) // generate/
	templatePath := filepath.Join(generateDir, "templates", templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %v", templateName, err)
	}

	return tmpl, nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func GeneratePrimaryInitScript(cfg *config.Config, outputDir string) error {
	// Load template
	tmpl, err := loadInitScriptTemplate("primary-init.sql.tmpl")
	if err != nil {
		return err
	}

	// Prepare template data
	data := InitScriptData{
		DbUser:              cfg.Primary.DbUser,
		DbPassword:          cfg.Primary.DbPassword,
		ReplicationUser:     cfg.Primary.ReplicationUser,
		ReplicationPassword: cfg.Primary.ReplicationPassword,
		DatadogPassword:     cfg.Monitoring.Datadog.DatadogUserPassword,
		WalLevel:            cfg.Options.WalLevel,
		MaxWalSenders:       cfg.Options.MaxWalSenders,
		WalKeepSize:         cfg.Options.WalKeepSize,
		HotStandby:          boolToString(cfg.Options.HotStandby),
		SynchronousCommit:   cfg.Options.SynchronousCommit,
	}

	outputFile := filepath.Join(outputDir, "01-setup-primary.sql")
	file, err := os.Create(outputFile)

	if err != nil {
		return fmt.Errorf("failed to create primary init script: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute primary init template: %v", err)
	}

	fmt.Printf("✅ Generated primary init script: %s\n", outputFile)
	return nil
}

func GenerateReplicaInitScripts(cfg *config.Config, outputDir string) error {
	tmpl, err := loadInitScriptTemplate("replica-init.sql.tmpl")
	if err != nil {
		return err
	}

	for i, replica := range cfg.Replicas {
		replicaName := fmt.Sprintf("replica%d", i+1)

		data := InitScriptData{
			DbUser:          replica.DbUser,
			DbPassword:      replica.DbPassword,
			DatadogPassword: cfg.Monitoring.Datadog.DatadogUserPassword,
			HotStandby:      boolToString(cfg.Options.HotStandby),
			ReplicationSlot: replica.ReplicationSlot,
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("01-setup-%s.sql", replicaName))
		file, err := os.Create(outputFile)

		if err != nil {
			return fmt.Errorf("failed to create %s init script: %v", replicaName, err)
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			return fmt.Errorf("failed to execute %s init template: %v", replicaName, err)
		}

		fmt.Printf("✅ Generated %s init script: %s\n", replicaName, outputFile)
	}

	return nil
}

func GenerateAllInitScripts(cfg *config.Config, outputDir string) error {
	fmt.Println("📝 Generating SQL initialization scripts...")

	if err := GeneratePrimaryInitScript(cfg, outputDir); err != nil {
		return fmt.Errorf("failed to generate primary init script: %v", err)
	}

	if err := GenerateReplicaInitScripts(cfg, outputDir); err != nil {
		return fmt.Errorf("failed to generate replica init scripts: %v", err)
	}

	fmt.Println("✅ All SQL initialization scripts generated successfully")
	return nil
}
