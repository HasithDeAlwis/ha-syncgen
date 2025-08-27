package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syncgen/devenv/aws/generate"
	"syncgen/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go generate-all-from-terraform <tf_output>    # Generate ALL files from terraform state (RECOMMENDED)")
		fmt.Println("  go run main.go generate-all-from-config                   # Generate ALL files from existing config")
		fmt.Println("  go run main.go generate-ssh-scripts                       # Generate SSH transfer and run scripts for VMs")
		fmt.Println("  go run main.go generate-syncgen-transfer-scripts          # Generate per-VM syncgen transfer/run scripts from generated/")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {

	case "generate-all-from-config":
		err := generateAllFromConfig()
		if err != nil {
			fmt.Printf("Error generating all files from config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("🎉 All files generated successfully - ready for deployment!")

	case "generate-all-from-terraform":
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run main.go generate-all-from-terraform <path_to_tf_output.json>")
			os.Exit(1)
		}

		err := generateAllFromTerraform(os.Args[2])
		if err != nil {
			fmt.Printf("Error generating all files from terraform: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("🎉 All files generated successfully - ready for deployment!")

	case "generate-syncgen-transfer-scripts":
		err := generateSyncgenTransferScripts()
		if err != nil {
			fmt.Printf("Error generating syncgen transfer scripts: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Syncgen transfer and run scripts generated successfully!")
	}
}

func generateSyncgenTransferScripts() error {
	fmt.Println("🔄 Generating syncgen transfer and run scripts from generated/ ...")

	// Load config from generated/config.yaml in root
	cfg, cfgErr := generate.LoadConfigFromGenerated()
	if cfgErr != nil {
		return fmt.Errorf("failed to load config from generated/config.yaml: %v", cfgErr)
	}
	localGenerate, genErr := localGenerateDirectory()
	if genErr != nil {
		return fmt.Errorf("failed to determine local generated/ directory: %v", genErr)
	}

	if err := generate.GenerateSyncgenTransferScripts(cfg, localGenerate); err != nil {
		return err
	}

	return nil
}

func runAllGenerators(cfg *config.Config, printStages bool) error {
	if printStages {
		fmt.Println("📝 Generating SQL initialization scripts...")
	}

	localGenerate, err := localGenerateDirectory()
	if err != nil {
		return fmt.Errorf("failed to determine local generated/ directory: %v", err)
	}

	if err := generate.GenerateAllInitScripts(cfg, localGenerate); err != nil {
		return fmt.Errorf("failed to generate SQL init scripts: %v", err)
	}

	if printStages {
		fmt.Println("🔧 Generating deployment scripts...")
	}
	if err := generate.GenerateDeploymentScripts(cfg); err != nil {
		return fmt.Errorf("failed to generate deployment scripts: %v", err)
	}

	return nil
}

func generateAllFromConfig() error {
	fmt.Println("🏗️  Generating all files from existing config...")

	// Load existing config
	cfg, err := generate.LoadConfigFromGenerated()
	if err != nil {
		return fmt.Errorf("failed to load config from generated/config.yaml: %v", err)
	}

	fmt.Printf("📋 Loaded config: primary %s with %d replicas\n",
		cfg.Primary.Host, len(cfg.Replicas))

	// Generate ALL files from in-memory config (prints stage headings)
	if err := runAllGenerators(cfg, true); err != nil {
		return err
	}

	fmt.Println("✅ All generation complete!")
	return nil
}

func generateAllFromTerraform(tfstatePath string) error {
	fmt.Println("🏗️  Generating all files from terraform state...")

	// Step 1: Parse terraform state directly into config struct (in memory)
	cfg, err := generate.ParseTFOutputsToConfig(tfstatePath)
	if err != nil {
		return fmt.Errorf("failed to parse terraform state: %v", err)
	}

	fmt.Printf("📋 Parsed terraform state: primary %s with %d replicas\n",
		cfg.Primary.Host, len(cfg.Replicas))

	// Step 2: Generate ALL files from in-memory config (prints stage headings)
	if err := runAllGenerators(cfg, true); err != nil {
		return err
	}

	// Step 3: Save config.yaml for syncgen compatibility (last step, not intermediate)
	fmt.Println("💾 Saving config.yaml for syncgen compatibility...")
	if err := generate.SaveConfigYAML(cfg); err != nil {
		return fmt.Errorf("failed to save config.yaml: %v", err)
	}

	fmt.Println("✅ All generation complete!")
	return nil
}

func localGenerateDirectory() (string, error) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}

	absPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return "", err
	}

	genDir := filepath.Join(filepath.Dir(absPath), "generated")
	genDirErr := os.MkdirAll(genDir, 0755)
	if genDirErr != nil {
		return "", fmt.Errorf("failed to create generated directory: %v", genDirErr)
	}
	return genDir, nil
}

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
	awsDir := filepath.Dir(absPath)        // ../aws
	devEnvDir := filepath.Dir(awsDir)      // ../devenv
	projectRoot := filepath.Dir(devEnvDir) // ../ha-syncgen
	genDir := filepath.Join(projectRoot, "generated")
	genDirErr := os.MkdirAll(genDir, 0755)
	if genDirErr != nil {
		return "", fmt.Errorf("failed to create generated directory: %v", err)
	}
	return genDir, nil
}
