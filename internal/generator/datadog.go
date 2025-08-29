package generator

import (
	"os"
	"path/filepath"
	"text/template"
)
func (g *Generator) generateDatadogFiles(installTmpl, sqlTmpl, confTmpl *template.Template) error {
	ddDir := filepath.Join(g.outputDir, "datadog")
	if err := os.MkdirAll(ddDir, 0755); err != nil {
		return err
	}
	if err := g.generateDDSQLFile(sqlTmpl, ddDir); err != nil {
		return err
	}
	return nil
}
// generateDDSQLFile creates the SQL file for Datadog user setup using a Go template
func (g *Generator) generateDDSQLFile(sqlTmpl *template.Template, ddDir string) error {
	data := map[string]interface{}{
		"Password": g.config.Primary.DbPassword,
	}
	outputFile := filepath.Join(ddDir, "datadog.sql")
	return executeTemplateToFile(sqlTmpl, data, outputFile, "datadog.sql")
}
