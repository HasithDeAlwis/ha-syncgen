package config

import (
	"fmt"
)

func Print(cfg *Config) {
	fmt.Printf("Build configuration:\n")
	fmt.Printf("Primary Cluster Node: %s:%d\n", cfg.Primary.Host, cfg.Primary.Port)
	if len(cfg.Replicas) == 0 {
		fmt.Println("No replicas configured.")
		return
	}
	for i, replica := range cfg.Replicas {
		fmt.Printf("\tReplica %d: %s (Sync Interval: %s)\n", i+1, replica.Host, replica.SyncInterval)
	}
}
