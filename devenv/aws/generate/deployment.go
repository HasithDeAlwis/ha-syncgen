package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"syncgen/internal/config"
)

// DeploymentData holds data for deployment script generation
type DeploymentData struct {
	PrimaryIP       string
	PrimaryUser     string
	PrimaryPassword string
	PrimaryDBName   string
	Replicas        []ReplicaDeploymentData
	SSHKeyPath      string
	SSHUser         string
}

type ReplicaDeploymentData struct {
	IP       string
	User     string
	Password string
	Name     string
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

func prepareDeploymentData(cfg *config.Config) DeploymentData {
	deployData := DeploymentData{
		PrimaryIP:       cfg.Primary.Host,
		PrimaryUser:     cfg.Primary.DbUser,
		PrimaryPassword: cfg.Primary.DbPassword,
		PrimaryDBName:   cfg.Primary.DbName,
		SSHKeyPath:      "$HOME/.ssh/your-aws-key.pem",
		SSHUser:         "ec2-user",
	}

	// Add replica data
	for i, replica := range cfg.Replicas {
		replicaData := ReplicaDeploymentData{
			IP:       replica.Host,
			User:     replica.DbUser,
			Password: replica.DbPassword,
			Name:     fmt.Sprintf("replica%d", i+1),
		}
		deployData.Replicas = append(deployData.Replicas, replicaData)
	}
	return deployData
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
	fmt.Println("🔧 Generating deployment scripts...")

	deployData := prepareDeploymentData(cfg)
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

	if err := generateAndWriteScript(outputDir, "DEPLOYMENT_README.md", "DEPLOYMENT_README.md.tmpl", deployData); err != nil {
		return err
	}

	fmt.Println("✅ All deployment scripts generated successfully")
	return nil
}

func generateScriptFromTemplate(tmplContent string, data DeploymentData, outputPath string) error {
	tmpl, err := template.New("script").Parse(tmplContent)
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
