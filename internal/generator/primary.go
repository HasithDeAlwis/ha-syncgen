package generator

import (
	"path/filepath"
	"text/template"
)

// generatePrimaryFiles creates configuration files for the primary PostgreSQL server
func (g *Generator) generatePrimaryFiles(pgHbaTmpl, postgresqlConfTmpl, setupPrimaryTmpl *template.Template) error {
	primaryDir := filepath.Join(g.outputDir, "primary")
	specs := []FileSpec{
		{
			Tmpl:     postgresqlConfTmpl,
			Dir:      primaryDir,
			Filename: "postgresql.conf.custom",
			Data: map[string]interface{}{
				"WalLevel":            g.config.Options.WalLevel,
				"HotStandby":          g.config.Options.HotStandby,
				"SynchronousCommit":   g.config.Options.SynchronousCommit,
				"MaxWalSenders":       g.config.Options.MaxWalSenders,
				"MaxReplicationSlots": len(g.config.Replicas) + 2,
				"WalKeepSize":         g.config.Options.WalKeepSize,
				"Port":                g.config.Primary.Port,
				"HasMonitoring":       g.config.Monitoring.Datadog.Enabled,
			},
		},
		{
			Tmpl:     pgHbaTmpl,
			Dir:      primaryDir,
			Filename: "pg_hba.conf.custom",
			Data: map[string]interface{}{
				"ReplicationUser": g.config.Primary.ReplicationUser,
				"Replicas":        g.config.Replicas,
			},
		},
		{
			Tmpl:     setupPrimaryTmpl,
			Dir:      primaryDir,
			Filename: "setup_primary.sh",
			Data: map[string]interface{}{
				"PrimaryHost":         g.config.Primary.Host,
				"PrimaryPort":         g.config.Primary.Port,
				"DbUser":              g.config.Primary.DbUser,
				"DbName":              g.config.Primary.DbName,
				"ReplicationUser":     g.config.Primary.ReplicationUser,
				"ReplicationPassword": g.config.Primary.ReplicationPassword,
				"Replicas":            g.config.Replicas,
				"DataDirectory":       g.config.Primary.DataDirectory,
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
