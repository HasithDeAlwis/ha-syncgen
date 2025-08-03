package generator

import (
	"fmt"
	"os"
	"path/filepath"

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
	if err := g.generatePrimaryFiles(); err != nil {
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
