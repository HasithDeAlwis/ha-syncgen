# WIP: ha-syncgen

> A declarative, YAML-powered tool for scaffolding high-availability PostgreSQL infrastructure using rsync and systemd.

## � Quick Start

1. **Install Go dependencies:**
```bash
go mod tidy
```

2. **Build the tool:**
```bash
go build -o syncgen .
```

3. **Create your cluster configuration:**
```yaml
# cluster.yaml
primary:
  host: 10.0.0.1
  port: 5432

replicas:
  - host: 10.0.0.2
    sync_interval: 30s
  - host: 10.0.0.3
    sync_interval: 60s

options:
  rsync_user: postgres
  sync_method: rsync
  promote_on_failure: true
  observability: datadog
```

4. **Generate HA infrastructure:**
```bash
./syncgen build cluster.yaml
```

This will create a `generated/` directory with all the necessary scripts and configuration files.

## 📁 Generated Output

```
generated/
├── replica-10.0.0.2/
│   ├── sync.sh              # Rsync synchronization script
│   ├── check_primary.sh     # Health check script  
│   ├── sync.service         # Systemd service unit
│   └── sync.timer           # Systemd timer unit
├── replica-10.0.0.3/
│   ├── sync.sh
│   ├── check_primary.sh
│   ├── sync.service
│   └── sync.timer
├── config/
│   ├── patch_postgresql_conf.sh  # PostgreSQL config patches
│   ├── patch_pg_hba_conf.sh      # pg_hba.conf patches
│   └── promote_replica.sh        # Manual promotion script
└── observability/
    ├── datadog-agent.yaml         # Datadog configuration
    ├── prometheus.yml             # Prometheus scrape config
    └── aggregate_status.sh        # JSON status aggregator
```

## 🛠️ Development

### Build for multiple platforms:
```bash
./scripts/build.sh
```

### Run with example configuration:
```bash
go run . build examples/cluster.yaml
```

## 📊 Features

- ✅ **YAML Configuration** - Simple, declarative cluster definition
- ✅ **Script Generation** - Creates all necessary sync and failover scripts
- ✅ **Systemd Integration** - Generates service and timer units
- ✅ **Health Monitoring** - Automatic primary failure detection
- ✅ **Observability** - Datadog, Prometheus, and JSON status support
- ✅ **Cross-platform** - Builds for Linux, macOS, and Windows

## 🎯 Roadmap

- [ ] `syncgen deploy` command with SSH automation
- [ ] `syncgen status` command for cluster monitoring  
- [ ] `syncgen promote` command for manual failover
- [ ] Advanced validation and error handling
- [ ] Integration tests with real PostgreSQL instances
- [ ] Docker support for development/testing
