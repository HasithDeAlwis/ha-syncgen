package generate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"syncgen/internal/config"

	"github.com/joho/godotenv"
)

// VMFileTransfer holds the mapping of files to transfer for a VM
type VMFileTransfer struct {
	Host      string
	Files     []string
	RemoteDir string
}

type SSHSyncgenTransferData struct {
	Transfers  []VMFileTransfer
	SSHKeyPath string
	SSHUser    string // Default to ec2-user if not set
}

type SSHRunScriptsData struct {
	SSHKeyPath string
	SSHUser    string
	Primary    VMFileTransfer
	Replicas   []VMFileTransfer
}

// PrettyPrint prints SSHSyncgenTransferData in a readable format
func (d *SSHSyncgenTransferData) PrettyPrint() {
	fmt.Printf("SSH User: %s\nSSH Key Path: %s\n", d.SSHUser, d.SSHKeyPath)
	fmt.Println("Transfers:")
	for _, t := range d.Transfers {
		fmt.Printf("  Host: %s\n", t.Host)
		fmt.Printf("  RemoteDir: %s\n", t.RemoteDir)
		fmt.Println("  Files:")
		for _, f := range t.Files {
			fmt.Printf("    - %s\n", f)
		}
	}
}

// getRootGeneratedDir returns the absolute path to the root-level generated/ directory
func getRootGeneratedDir() (string, error) {
	// Find the root of the repo by walking up from this file
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}
	absPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return "", err
	}
	// go up to repo root
	generateDir := filepath.Dir(absPath)   // ../generate
	awsDir := filepath.Dir(generateDir)    // ../aws
	devEnvDir := filepath.Dir(awsDir)      // ../devenv
	projectRoot := filepath.Dir(devEnvDir) // ../ha-syncgen
	genDir := filepath.Join(projectRoot, "generated")
	return genDir, nil
}
func collectAndPrepend(genDir, name string) ([]string, error) {
	dir := filepath.Join(genDir, name)
	files, err := listAllFiles(dir)
	if err != nil {
		return nil, err
	}
	return prependDir(name, files), nil
}

func collectVMFileTransfers(cfg *config.Config) ([]VMFileTransfer, error) {
	genDir, err := getRootGeneratedDir()
	if err != nil {
		return nil, err
	}
	var transfers []VMFileTransfer

	primaryFiles, err := collectAndPrepend(genDir, "primary")
	if err != nil {
		return nil, err
	}
	datadogFiles, err := collectAndPrepend(genDir, "datadog")
	if err != nil {
		return nil, err
	}
	allPrimaryFiles := append(primaryFiles, datadogFiles...)

	transfers = append(transfers, VMFileTransfer{
		Host:      cfg.Primary.Host,
		Files:     allPrimaryFiles,
		RemoteDir: "/syncgen",
	})

	for _, replica := range cfg.Replicas {
		replicaFiles, err := collectAndPrepend(genDir, "replica-"+replica.Host)
		if err != nil {
			return nil, err
		}
		allReplicaFiles := append(replicaFiles, datadogFiles...)
		transfers = append(transfers, VMFileTransfer{
			Host:      replica.Host,
			Files:     allReplicaFiles,
			RemoteDir: "/syncgen",
		})
	}

	return transfers, nil
}

// listAllFiles returns all file names (not directories) in a directory (relative to that dir)
func listAllFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(dir, path)
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}

func prependDir(dir string, files []string) []string {
	var out []string
	for _, f := range files {
		out = append(out, filepath.Join(dir, f))
	}
	return out
}

// GenerateSyncgenTransferScripts generates transfer-ha-scripts.sh and run-ha-scripts.sh for syncgen
func GenerateSyncgenTransferScripts(cfg *config.Config, genDir string) error {
	data, err := prepareSSHScriptData(cfg)

	sshKeyPath, okKey := os.LookupEnv("SSH_KEY_PATH")
	if !okKey || sshKeyPath == "" {
		return fmt.Errorf("SSH_KEY_PATH variable is not set in your .env file")
	}

	if err != nil {
		return err
	}

	tmplDir := filepath.Join(filepath.Dir(genDir), "generate", "templates")

	transferTmplPath := filepath.Join(tmplDir, "transfer-ha-scripts.sh.tmpl")
	runTmplPath := filepath.Join(tmplDir, "run-ha-scripts.sh.tmpl")

	transferPath := filepath.Join(genDir, "transfer-ha-scripts.sh")
	runPath := filepath.Join(genDir, "run-ha-scripts.sh")

	if err := renderSyncgenScriptTemplate(transferTmplPath, data, transferPath); err != nil {
		return fmt.Errorf("failed to generate transfer-ha-scripts.sh: %w", err)
	}

	modifiedData := map[string]interface{}{
		"SSHKeyPath": data.SSHKeyPath,
		"SSHUser":    data.SSHUser,
		"Primary":    data.Transfers[0],
		"Replicas":   data.Transfers[1:],
	}
	if err := renderSyncgenScriptTemplate(runTmplPath, modifiedData, runPath); err != nil {
		return fmt.Errorf("failed to generate run-ha-scripts.sh: %w", err)
	}

	fmt.Printf("✅ Generated syncgen SSH scripts: %s, %s\n", transferPath, runPath)
	return nil
}

func renderSyncgenScriptTemplate(tmplPath string, data interface{}, outPath string) error {
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}
	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplPath, err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to render template %s: %w", tmplPath, err)
	}

	if err := os.Chmod(outPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", outPath, err)
	}
	return nil
}

// prepareSSHScriptData builds SSHScriptData from config and env
func prepareSSHScriptData(cfg *config.Config) (*SSHSyncgenTransferData, error) {
	_ = godotenv.Load("../../.env")

	sshUser, okUser := os.LookupEnv("SSH_USER")
	sshKeyPath, okKey := os.LookupEnv("SSH_KEY_PATH")

	if !okKey || sshKeyPath == "" {
		return nil, fmt.Errorf("SSH_KEY_PATH variable is not set in your .env file")
	}

	if !okUser || sshUser == "" {
		sshUser = "ec2-user"
	}

	transfers, err := collectVMFileTransfers(cfg)

	if err != nil {
		return nil, err
	}

	data := SSHSyncgenTransferData{
		Transfers:  transfers,
		SSHKeyPath: sshKeyPath,
		SSHUser:    sshUser,
	}
	return &data, nil
}
