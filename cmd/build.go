package cmd

import (
	"fmt"
	"os"
	"syncgen/internal/config"
	"syncgen/internal/generator"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [config file]",
	Short: "Generate HA scripts and configuration from cluster.yaml",
	Long: `Build reads your cluster.yaml configuration and generates all the necessary
files for high-availability PostgreSQL setup including:
- Sync scripts for each replica
- Systemd service and timer units  
- Health check and failover scripts
- Optional observability configurations`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,

	Run: func(cmd *cobra.Command, args []string) {
		configFile := args[0]
		cfg, err := config.Parse(configFile)
		if err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			return
		}

		// Print the configuration for user verification
		fmt.Println("Configuration validated successfully:")
		config.Print(cfg)

		// Generate output directory
		outputDir := "generated"
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			return
		}

		// Initialize generator and create all files
		gen := generator.New(cfg, outputDir)
		if err := gen.GenerateAll(); err != nil {
			fmt.Printf("Error generating files: %v\n", err)
			return
		}

		fmt.Printf("\n✅ HA PostgreSQL configuration generated successfully in '%s/' directory\n", outputDir)
		fmt.Println("\nNext steps:")
		fmt.Printf("1. Review the generated scripts in '%s/'\n", outputDir)
		fmt.Println("2. Copy the scripts to your target servers")
		fmt.Println("3. Set up PostgreSQL streaming replication by running the setup scripts")
		fmt.Println("4. Install and enable the systemd services for automatic health monitoring")
		fmt.Println("\nFor deployment help, run: syncgen --help")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	// No flags configured for build command yet
}
