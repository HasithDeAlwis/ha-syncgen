package generator

import (
	"fmt"
	"path/filepath"
	"syncgen/internal/config"
)

// generateSystemdService creates a systemd service unit for PostgreSQL replication management
func (g *Generator) generateSystemdService(replica config.Replica, replicaDir string) error {
	service := fmt.Sprintf(`[Unit]
Description=PostgreSQL HA Health Check for replica %s
Documentation=https://github.com/HasithDeAlwis/ha-syncgen
After=postgresql.service network.target
Requires=postgresql.service
PartOf=postgresql.service

[Service]
Type=oneshot
User=postgres
Group=postgres
ExecStart=%s/health_check.sh
WorkingDirectory=%s
StandardOutput=journal
StandardError=journal
SyslogIdentifier=ha-syncgen-health-check

# Environment variables for the health check script
Environment=PGUSER=%s
Environment=PGDATABASE=postgres
Environment=REPLICA_HOST=%s
Environment=PRIMARY_HOST=%s

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/ha-syncgen %s

# Restart behavior
RemainAfterExit=no
TimeoutStartSec=60
TimeoutStopSec=30

[Install]
WantedBy=multi-user.target
`,
		replica.Host,
		replicaDir,
		replicaDir,
		g.config.Primary.ReplicationUser,
		replica.Host,
		g.config.Primary.Host,
		g.config.Primary.DataDirectory)

	return g.writeFile(filepath.Join(replicaDir, "ha-postgres-health.service"), service)
}

// generateSystemdTimer creates a systemd timer unit for regular health checks
func (g *Generator) generateSystemdTimer(replica config.Replica, replicaDir string) error {
	timer := fmt.Sprintf(`[Unit]
Description=PostgreSQL HA Health Check Timer for replica %s
Documentation=https://github.com/HasithDeAlwis/ha-syncgen
Requires=ha-postgres-health.service
After=postgresql.service

[Timer]
# Run health check every 30 seconds
OnBootSec=2min
OnUnitActiveSec=30sec

# Ensure timer is persistent across reboots
Persistent=true

# Prevent timer drift
AccuracySec=5sec

# Randomize execution slightly to prevent thundering herd
RandomizedDelaySec=5sec

[Install]
WantedBy=timers.target
`,
		replica.Host)

	return g.writeFile(filepath.Join(replicaDir, "ha-postgres-health.timer"), timer)
}
