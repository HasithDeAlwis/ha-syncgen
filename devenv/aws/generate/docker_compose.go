package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"syncgen/internal/config"
)

// TemplateData holds the data for Docker Compose template rendering
type DockerComposeData struct {
	DbName            string
	DbUser            string
	DbPassword        string
	WalLevel          string
	MaxWalSenders     int
	WalKeepSize       string
	HotStandby        string
	SynchronousCommit string
}

// generateDockerComposeDir creates the full directory path for a given server type
func generateDockerComposeDir(serverType string) (string, error) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}

	absPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return "", err
	}

	generateDir := filepath.Dir(absPath) // generate/
	awsDir := filepath.Dir(generateDir)  // aws/
	generatedDir := filepath.Join(awsDir, "generated", serverType)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Join(generatedDir, "init-scripts"), 0755); err != nil {
		return "", err
	}

	return generatedDir, nil
}

// loadDockerComposeTemplate loads a template file from the templates directory
func loadDockerComposeTemplate(templateName string) (*template.Template, error) {
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
func writeDockerComposeFile(tmpl *template.Template, cfg *config.Config, file *os.File) error {
	data := DockerComposeData{
		DbName:            cfg.Primary.DbName,
		DbUser:            cfg.Primary.DbUser,
		DbPassword:        cfg.Primary.DbPassword,
		WalLevel:          cfg.Options.WalLevel,
		MaxWalSenders:     cfg.Options.MaxWalSenders,
		WalKeepSize:       cfg.Options.WalKeepSize,
		HotStandby:        boolToString(cfg.Options.HotStandby),
		SynchronousCommit: cfg.Options.SynchronousCommit,
	}
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	return nil
}

func GeneratePrimaryDockerCompose(cfg *config.Config) error {
	// Load template
	tmpl, err := loadDockerComposeTemplate("docker-compose.primary.yml.tmpl")
	if err != nil {
		return err
	}

	// Create output directory
	outputDir, err := generateDockerComposeDir("primary")
	if err != nil {
		return err
	}
	outputFile := filepath.Join(outputDir, "docker-compose.yml")
	file, err := os.Create(outputFile)
	defer file.Close()

	// Prepare template data
	if err := writeDockerComposeFile(tmpl, cfg, file); err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to create primary docker-compose file: %v", err)
	}

	fmt.Printf("✅ Generated primary docker-compose.yml: %s\n", outputFile)
	return nil
}

// GenerateReplicaDockerCompose generates Docker Compose files for replica databases
func GenerateReplicaDockerCompose(cfg *config.Config) error {
	// Load template
	tmpl, err := loadDockerComposeTemplate("docker-compose.replica.yml.tmpl")
	if err != nil {
		return err
	}

	// Generate for each replica
	for i, replica := range cfg.Replicas {
		replicaName := fmt.Sprintf("replica%d", i+1)

		replicaCfg := *cfg // shallow copy
		replicaCfg.Primary.DbUser = replica.DbUser
		replicaCfg.Primary.DbPassword = replica.DbPassword

		// Create output directory
		outputDir, err := generateDockerComposeDir(replicaName)
		if err != nil {
			return err
		}

		// Generate docker-compose.yml
		outputFile := filepath.Join(outputDir, "docker-compose.yml")
		file, err := os.Create(outputFile)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("failed to create replica docker-compose file: %v", err)
		}

		if err := writeDockerComposeFile(tmpl, &replicaCfg, file); err != nil {
			return fmt.Errorf("failed to execute replica template: %v", err)
		}

		fmt.Printf("✅ Generated %s docker-compose.yml: %s\n", replicaName, outputFile)
	}

	return nil
}

// GenerateAllDockerCompose generates all Docker Compose files
func GenerateAllDockerCompose(cfg *config.Config) error {
	fmt.Println("🐳 Generating Docker Compose files...")

	if err := GeneratePrimaryDockerCompose(cfg); err != nil {
		return fmt.Errorf("failed to generate primary docker-compose: %v", err)
	}

	if err := GenerateReplicaDockerCompose(cfg); err != nil {
		return fmt.Errorf("failed to generate replica docker-compose files: %v", err)
	}

	fmt.Println("✅ All Docker Compose files generated successfully")
	return nil
}

// Helper function to convert bool to string for PostgreSQL config
func boolToString(b bool) string {
	if b {
		return "on"
	}
	return "off"
}
