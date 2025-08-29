package generator

import (
	"os"
	"path/filepath"
	"text/template"
)

// generatePrimaryFiles creates configuration files for the primary PostgreSQL server
func (g *Generator) generatePrimaryFiles(pgHbaTmpl, postgresqlConfTmpl, setupPrimaryTmpl *template.Template) error {
	primaryDir := filepath.Join(g.outputDir, "primary")
	if err := os.MkdirAll(primaryDir, 0755); err != nil {
		return err
	}

	// Generate PostgreSQL configuration patches
	if err := g.generatePostgreSQLConfig(primaryDir, postgresqlConfTmpl); err != nil {
		return err
	}

	// Generate pg_hba.conf patches for replication
	if err := g.generatePgHbaConfig(primaryDir, pgHbaTmpl); err != nil {
		return err
	}

	// Generate replication slot setup script
	if err := g.generateReplicationSlotSetup(primaryDir, setupPrimaryTmpl); err != nil {
		return err
	}

	return nil
}

// generatePostgreSQLConfig creates postgresql.conf patches for streaming replication

func (g *Generator) generatePostgreSQLConfig(primaryDir string, postgresqlConfTmpl *template.Template) error {
	data := map[string]interface{}{
		"WalLevel":            g.config.Options.WalLevel,
		"HotStandby":          g.config.Options.HotStandby,
		"SynchronousCommit":   g.config.Options.SynchronousCommit,
		"MaxWalSenders":       g.config.Options.MaxWalSenders,
		"MaxReplicationSlots": len(g.config.Replicas) + 2,
		"WalKeepSize":         g.config.Options.WalKeepSize,
		"Port":                g.config.Primary.Port,
		"HasMonitoring":       g.config.Monitoring.Datadog.Enabled,
	}
	outputFile := filepath.Join(primaryDir, "postgresql.conf.custom")
	return executeTemplateToFile(postgresqlConfTmpl, data, outputFile, "postgres.conf")
}

// generatePgHbaConfig creates pg_hba.conf entries for replication using a Go template

func (g *Generator) generatePgHbaConfig(primaryDir string, pgHbaTmpl *template.Template) error {
	data := map[string]interface{}{
		"ReplicationUser": g.config.Primary.ReplicationUser,
		"Replicas":        g.config.Replicas,
	}
	outputFile := filepath.Join(primaryDir, "pg_hba.conf.custom")
	return executeTemplateToFile(pgHbaTmpl, data, outputFile, "pg_hba.conf")
}

// generateReplicationSlotSetup creates a script to set up replication slots on the primary using a Go template

func (g *Generator) generateReplicationSlotSetup(primaryDir string, setupPrimaryTmpl *template.Template) error {
	data := map[string]interface{}{
		"PrimaryHost":         g.config.Primary.Host,
		"PrimaryPort":         g.config.Primary.Port,
		"DbUser":              g.config.Primary.DbUser,
		"DbName":              g.config.Primary.DbName,
		"ReplicationUser":     g.config.Primary.ReplicationUser,
		"ReplicationPassword": g.config.Primary.ReplicationPassword,
		"Replicas":            g.config.Replicas,
		"DataDirectory":       g.config.Primary.DataDirectory,
	}
	outputFile := filepath.Join(primaryDir, "setup_primary.sh")
	return executeTemplateToFile(setupPrimaryTmpl, data, outputFile, "setup_primary.sh")
}
