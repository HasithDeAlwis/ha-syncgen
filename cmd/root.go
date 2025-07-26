/*
Copyright © 2025 ha-syncgen contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syncgen",
	Short: "PostgreSQL High Availability synchronization and failover automation tool",
	Long: `syncgen is a command-line tool for setting up and managing PostgreSQL high availability
configurations using file-level synchronization and automated failover mechanisms.

The tool generates the necessary scripts, configuration files, and systemd services to
establish a primary-replica PostgreSQL setup with:
- Automated file synchronization using rsync over SSH
- Real-time monitoring and health checks
- Automatic failover to replica on primary failure
- SSH key-based authentication for secure data transfer
- Configurable sync intervals and retry policies

Usage examples:
  syncgen validate ha-config.yaml # Validate configuration
  syncgen build [config file]               # Generate HA scripts and configuration
  //  COMING SOON
  syncgen status                           # Check current HA status
  syncgen failover                         # Manual failover to replica

For more information, visit: https://github.com/HasithDeAlwis/ha-syncgen`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global configuration flags that apply to all subcommands
	// These flags will be available across all ha-syncgen commands

	// Config file flag - allows users to specify custom config location
	// Default behavior will look for config in a standard location:
	// 1. ./cluster.yaml (current directory)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default searches standard locations)")

	// Verbose output flag for debugging and detailed operation logs
	// rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output for debugging")

	// Dry-run flag to preview operations without making actual changes
	// rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview operations without executing them")
}
