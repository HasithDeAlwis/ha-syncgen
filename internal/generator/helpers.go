package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// executeTemplateToFile is a helper to execute a template with data and write to a file, with error context
func executeTemplateToFile(tmpl *template.Template, data interface{}, outputPath string, templateName string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %v", templateName, err)
	}

	// Make scripts executable if .sh
	if filepath.Ext(outputPath) == ".sh" {
		if err := os.Chmod(outputPath, 0755); err != nil {
			return fmt.Errorf("failed to make script executable: %w", err)
		}
	}

	return nil
}

// parseTemplateByName loads and parses a template by name from the given directory
func parseTemplateByName(templateDirectory, templateName string) (*template.Template, error) {
	tmplPath := filepath.Join(templateDirectory, templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s template: %w", templateName, err)
	}
	return tmpl, nil
}
