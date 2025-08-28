package generator

import (
	"path/filepath"
	"syncgen/internal/config"
)

// generateSystemdService creates a systemd service unit for PostgreSQL replication management using a template
func (g *Generator) generateSystemdService(replica config.Replica, replicaDir string) error {
	tmplDir, tmplErr := getTemplateDirectory()
	if tmplErr != nil {
		return tmplErr
	}
	serviceTmpl, err := parseTemplateByName(tmplDir, "ha-postgres-health.service.tmpl")
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Replica":    replica,
		"Primary":    g.config.Primary,
		"ReplicaDir": replicaDir,
	}
	outputFile := filepath.Join(replicaDir, "ha-postgres-health.service")
	return executeTemplateToFile(serviceTmpl, data, outputFile, "ha-postgres-health.service")
}

// generateSystemdTimer creates a systemd timer unit for regular health checks using a template
func (g *Generator) generateSystemdTimer(replica config.Replica, replicaDir string) error {
	tmplDir, tmplErr := getTemplateDirectory()
	if tmplErr != nil {
		return tmplErr
	}
	timerTmpl, err := parseTemplateByName(tmplDir, "ha-postgres-health.timer.tmpl")
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Replica": replica,
	}
	outputFile := filepath.Join(replicaDir, "ha-postgres-health.timer")
	return executeTemplateToFile(timerTmpl, data, outputFile, "ha-postgres-health.timer")
}
