package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/joho/godotenv"

	"syncgen/internal/config"
)

// DeploymentData holds data for deployment script generation
type DeploymentData struct {
	PrimaryIP  string
	Replicas   []string
	SSHKeyPath string
	SSHUser    string
}

func loadDeploymentTemplate(templateName string) (string, error) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}
	absPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return "", err
	}
	generateDir := filepath.Dir(absPath) // generate/
	templatePath := filepath.Join(generateDir, "templates", templateName)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %v", templateName, err)
	}
	return string(content), nil
}

func generateDeploymentDir() (string, error) {
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
	generatedDir := filepath.Join(awsDir, "generated")

	// Ensure directory exists
	if err := os.MkdirAll(generatedDir, 0755); err != nil {
		return "", err
	}

	return generatedDir, nil
}

func prepareDeploymentData(cfg *config.Config) (DeploymentData, error) {
	// godotenv.Load() should have already been called in GenerateDeploymentScripts

	sshUser, okUser := os.LookupEnv("SSH_USER")
	sshKeyPath, okKey := os.LookupEnv("SSH_KEY_PATH")

	if !okKey || sshKeyPath == "" {
		return DeploymentData{}, fmt.Errorf("ssh_key_path variable is not set in your .env file")
	}

	if !okUser || sshUser == "" {
		sshUser = "ec2-user"
	}

	deployData := DeploymentData{
		PrimaryIP:  cfg.Primary.Host,
		SSHKeyPath: sshKeyPath,
		SSHUser:    sshUser,
	}

	for _, replica := range cfg.Replicas {
		replicaData := replica.Host
		deployData.Replicas = append(deployData.Replicas, replicaData)
	}
	return deployData, nil
}

func generateAndWriteScript(outputDir string, outputFileName string, templateName string, data DeploymentData) error {
	templateContent, err := loadDeploymentTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %v", templateName, err)
	}

	outputPath := filepath.Join(outputDir, outputFileName)
	if err := generateScriptFromTemplate(templateContent, data, outputPath); err != nil {
		return fmt.Errorf("failed to generate script %s: %v", outputPath, err)
	}
	if err := os.Chmod(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}
	fmt.Printf("✅ Generated script: %s\n", outputPath)
	return nil
}

func GenerateDeploymentScripts(cfg *config.Config) error {
	// Load .env file if present
	_ = godotenv.Load("../../.env")

	fmt.Println("🔧 Generating deployment scripts...")

	deployData, err := prepareDeploymentData(cfg)
	if err != nil {
		return err
	}
	// Get output directory
	outputDir, err := generateDeploymentDir()
	if err != nil {
		return err
	}

	if err := generateAndWriteScript(outputDir, "deploy-to-servers.sh", "deploy-to-servers.sh.tmpl", deployData); err != nil {
		return err
	}

	if err := generateAndWriteScript(outputDir, "deploy-databases.sh", "deploy-databases.sh.tmpl", deployData); err != nil {
		return err
	}

	fmt.Println("✅ All deployment scripts generated successfully")
	return nil
}

func generateScriptFromTemplate(tmplContent string, data DeploymentData, outputPath string) error {
	funcMap := template.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}
	tmpl, err := template.New("script").Funcs(funcMap).Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create script file: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
