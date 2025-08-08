package config

import (
	"fmt"
)

func Print(cfg *Config) {
	fmt.Printf("=== PostgreSQL HA Streaming Replication Configuration ===\n\n")

	// Print primary configuration
	fmt.Printf("Primary Server:\n")
	fmt.Printf("  Host: %s:%d\n", cfg.Primary.Host, cfg.Primary.Port)
	fmt.Printf("  Data Directory: %s\n", cfg.Primary.DataDirectory)
	fmt.Printf("  Replication User: %s\n", cfg.Primary.ReplicationUser)
	fmt.Printf("  Password: %s\n", maskPassword(cfg.Primary.ReplicationPassword))

	// Print replica configuration
	if len(cfg.Replicas) == 0 {
		fmt.Printf("\nReplicas: None configured\n")
	} else {
		fmt.Printf("\nReplicas (%d configured):\n", len(cfg.Replicas))
		for i, replica := range cfg.Replicas {
			fmt.Printf("  %d. %s:%d\n", i+1, replica.Host, replica.Port)
			fmt.Printf("     Replication Slot: %s\n", replica.ReplicationSlot)
			fmt.Printf("     Sync Mode: %s\n", replica.SyncMode)
		}
	}

	// Print options configuration
	fmt.Printf("\nPostgreSQL Streaming Options:\n")
	fmt.Printf("  WAL Level: %s\n", cfg.Options.WalLevel)
	fmt.Printf("  Max WAL Senders: %d\n", cfg.Options.MaxWalSenders)
	fmt.Printf("  WAL Keep Size: %s\n", cfg.Options.WalKeepSize)
	fmt.Printf("  Hot Standby: %t\n", cfg.Options.HotStandby)
	fmt.Printf("  Synchronous Commit: %s\n", cfg.Options.SynchronousCommit)
	fmt.Printf("  Auto-promote on Primary Failure: %t\n", cfg.Options.PromoteOnFailure)

	fmt.Printf("\n=== Configuration Summary ===\n")
	fmt.Printf("Total nodes: %d (1 primary + %d replicas)\n", 1+len(cfg.Replicas), len(cfg.Replicas))
	fmt.Printf("Replication type: PostgreSQL Streaming Replication\n")
	if cfg.Options.PromoteOnFailure {
		fmt.Printf("Failover: Automatic promotion enabled\n")
	} else {
		fmt.Printf("Failover: Manual promotion only\n")
	}
}

// maskPassword partially hides the password for security
func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
