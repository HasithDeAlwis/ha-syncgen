package generator

import (
	"path/filepath"
	"text/template"
)

func (g *Generator) generateDatadogFiles(installTmpl, sqlTmpl, confTmpl *template.Template) error {
	ddDir := filepath.Join(g.outputDir, "datadog")
	specs := []FileSpec{
		{
			Tmpl:     installTmpl,
			Dir:      ddDir,
			Filename: "datadog-install.sh",
			Data: map[string]interface{}{
				"DataDogApiKey": g.config.Monitoring.Datadog.ApiKey,
				"DataDogSite":   g.config.Monitoring.Datadog.Site,
			},
		},
		{
			Tmpl:     sqlTmpl,
			Dir:      ddDir,
			Filename: "datadog.sql",
			Data: map[string]interface{}{
				"Password": g.config.Primary.DbPassword,
			},
		},
		{
			Tmpl:     confTmpl,
			Dir:      ddDir,
			Filename: "datadog-conf.yaml",
			Data: map[string]interface{}{
				"Host":          g.config.Primary.Host,
				"Port":          g.config.Primary.Port,
				"DbName":        g.config.Primary.DbName,
				"DbUser":        g.config.Primary.DbUser,
				"DbPassword":    g.config.Primary.DbPassword,
				"Datadirectory": g.config.Primary.DataDirectory,
			},
		},
	}
	for _, spec := range specs {
		if err := g.renderTemplate(spec.Tmpl, spec.Dir, spec.Filename, spec.Data); err != nil {
			return err
		}
	}
	return nil
}
