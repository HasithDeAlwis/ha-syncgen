package generator

import (
	"path/filepath"
	"syncgen/internal/config"
)

// generateSyncScript creates the streaming replication setup script using a Go template
func (g *Generator) generateSyncScript(replica config.Replica, replicaDir string) error {
	tmplDir, tmplErr := getTemplateDirectory()
	if tmplErr != nil {
		return tmplErr
	}
	syncScriptTmpl, err := parseTemplateByName(tmplDir, "setup_replication.sh.tmpl")
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Replica": replica,
		"Primary": g.config.Primary,
	}
	outputFile := filepath.Join(replicaDir, "setup_replication.sh")
	return executeTemplateToFile(syncScriptTmpl, data, outputFile, "setup_replication.sh")
}
