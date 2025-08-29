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
	if err := g.generateDDInstallScript(installTmpl, ddDir); err != nil {
		return err
	}
	if err := g.generateDDSQLFile(sqlTmpl, ddDir); err != nil {
		return err
	}
	if err := g.generateDDConfigFile(confTmpl, ddDir); err != nil {
		return err
	}
	return nil
}

// generateDDInstallScript creates the Datadog agent install script using a Go template
func (g *Generator) generateDDInstallScript(installTmpl *template.Template, ddDir string) error {
	data := map[string]interface{}{
		"DataDogApiKey": g.config.Monitoring.Datadog.ApiKey,
		"DataDogSite":   g.config.Monitoring.Datadog.Site,
	}
	outputFile := filepath.Join(ddDir, "datadog-install.sh")
	return executeTemplateToFile(installTmpl, data, outputFile, "datadog-install.sh")
}

// generateDDSQLFile creates the SQL file for Datadog user setup using a Go template
func (g *Generator) generateDDSQLFile(sqlTmpl *template.Template, ddDir string) error {
	data := map[string]interface{}{
		"Password": g.config.Primary.DbPassword,
	}
	outputFile := filepath.Join(ddDir, "datadog.sql")
	return executeTemplateToFile(sqlTmpl, data, outputFile, "datadog.sql")
}

// generateDDConfigFile creates the Datadog agent config YAML using a Go template
func (g *Generator) generateDDConfigFile(confTmpl *template.Template, ddDir string) error {
	data := map[string]interface{}{
		"Host":          g.config.Primary.Host,
		"Port":          g.config.Primary.Port,
		"DbName":        g.config.Primary.DbName,
		"DbUser":        g.config.Primary.DbUser,
		"DbPassword":    g.config.Primary.DbPassword,
		"Datadirectory": g.config.Primary.DataDirectory,
	}
	outputFile := filepath.Join(ddDir, "datadog-conf.yaml")
	return executeTemplateToFile(confTmpl, data, outputFile, "datadog-conf.yaml")
}
