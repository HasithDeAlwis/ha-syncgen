package generator

import (
	"path/filepath"
	"syncgen/internal/config"
)

// generateHealthCheckScript creates a health check script for monitoring primary PostgreSQL status using a template
func (g *Generator) generateHealthCheckScript(replica config.Replica, replicaDir string) error {
	tmplDir, tmplErr := getTemplateDirectory()
	if tmplErr != nil {
		return tmplErr
	}
	healthTmpl, err := parseTemplateByName(tmplDir, "health_check.sh.tmpl")
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Replica": replica,
		"Primary": g.config.Primary,
		"Options": g.config.Options,
	}
	outputFile := filepath.Join(replicaDir, "health_check.sh")
	return executeTemplateToFile(healthTmpl, data, outputFile, "health_check.sh")
}
