package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"syncgen/internal/config"
)

// Generator handles the generation of all configuration files
type Generator struct {
	config    *config.Config
	outputDir string
}

// New creates a new Generator instance
func New(cfg *config.Config, outputDir string) *Generator {
	return &Generator{
		config:    cfg,
		outputDir: outputDir,
	}
}

// GenerateAll generates all necessary files for the HA setup
func (g *Generator) GenerateAll() error {
	// Generate primary configuration files
	tmplDir, tmplErr := getTemplateDirectory()
	if tmplErr != nil {
		return fmt.Errorf("failed to get template directory: %w", tmplErr)
	}

	pgHbaTmpl, err := parseTemplateByName(tmplDir, "pg_hba.conf.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse pg_hba.conf template: %w", err)
	}

	postgresqlConfTmpl, err := parseTemplateByName(tmplDir, "postgresql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse postgresql.conf template: %w", err)
	}

	setupPrimaryTmpl, err := parseTemplateByName(tmplDir, "setup_primary.sh.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse setup_primary.sh template: %w", err)
	}

	if err := g.generatePrimaryFiles(pgHbaTmpl, postgresqlConfTmpl, setupPrimaryTmpl); err != nil {
		return fmt.Errorf("failed to generate primary files: %w", err)
	}

	// Generate replica-specific files
	for _, replica := range g.config.Replicas {
		if err := g.generateReplicaFiles(replica); err != nil {
			return fmt.Errorf("failed to generate files for replica %s: %w", replica.Host, err)
		}
	}

	return nil
}

// generateReplicaFiles generates all files specific to a replica
func (g *Generator) generateReplicaFiles(replica config.Replica) error {
	replicaDir := filepath.Join(g.outputDir, fmt.Sprintf("replica-%s", replica.Host))
	if err := os.MkdirAll(replicaDir, 0755); err != nil {
		return err
	}

	// Generate sync script
	if err := g.generateSyncScript(replica, replicaDir); err != nil {
		return err
	}

	// Generate health check script
	if err := g.generateHealthCheckScript(replica, replicaDir); err != nil {
		return err
	}

	// Generate systemd service
	if err := g.generateSystemdService(replica, replicaDir); err != nil {
		return err
	}

	// Generate systemd timer
	if err := g.generateSystemdTimer(replica, replicaDir); err != nil {
		return err
	}

	return nil
}

func getTemplateDirectory() (string, error) {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}
	absPath, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}
	generatorDir := filepath.Dir(absPath)
	return filepath.Join(generatorDir, "templates"), nil
}

func getTemplate(templateDirectory string, templateName string) (string, error) {
	templatePath := filepath.Join(templateDirectory, templateName)
	return templatePath, nil
}
