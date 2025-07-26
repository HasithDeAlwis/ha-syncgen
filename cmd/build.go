package cmd

import (
	"fmt"
	"ha-syncgen/internal/config"

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
		config.Print(cfg)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	// No flags configured for build command yet
}
