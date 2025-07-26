/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"syncgen/internal/config"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file without generating files",
	Long: `The validate command checks the provided configuration file for correctness.
	
	Example usage:
# cluster.yaml
primary:
  host: 10.0.0.1
  port: 5432

replicas:
  - host: 10.0.0.2
    sync_interval: 30s
  - host: 10.0.0.3
    sync_interval: 60s
	syncgen validate cluster.yaml

> syncgen validate cluster.yaml
Build configuration:
  Cluster Name: 10.0.0.1
  Replica: 10.0.0.2 (Sync Interval: 30s)
  Replica: 10.0.0.3 (Sync Interval: 60s)

# cluster.yaml
primary: 
  port: 5432

> syncgen validate cluster.yaml
Error validating config file:
Error: Primary host is required in the configuration file. Replicas are not provided.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := args[0]
		cfg, err := config.Parse(configFile)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("Validation successful!")
		config.Print(cfg)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
