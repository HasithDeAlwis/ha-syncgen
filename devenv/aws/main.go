package main

import (
	"fmt"
	"os"
	"syncgen/devenv/aws/generate"
	"syncgen/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go generate-all-from-terraform <tf_output>    # Generate ALL files from terraform state (RECOMMENDED)")
		fmt.Println("  go run main.go generate-all-from-config                   # Generate ALL files from existing config")
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
	}
}

func runAllGenerators(cfg *config.Config, printStages bool) error {
	if printStages {
		fmt.Println("🐳 Generating Docker Compose files...")
	}
	if err := generate.GenerateAllDockerCompose(cfg); err != nil {
		return fmt.Errorf("failed to generate Docker Compose files: %v", err)
	}

	if printStages {
		fmt.Println("📝 Generating SQL initialization scripts...")
	}
	if err := generate.GenerateAllInitScripts(cfg); err != nil {
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
