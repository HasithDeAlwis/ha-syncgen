
# ha-syncgen

PostgreSQL High Availability synchronization and failover automation tool

---

## Overview

`ha-syncgen` is a CLI tool for declaratively managing PostgreSQL high availability (HA) clusters. It generates scripts, configuration files, and systemd services to automate primary-replica setups, health checks, and failover. The tool is powered by a simple YAML configuration and is designed for reliability, transparency, and easy integration with your infrastructure workflows.

**Key Features:**
- **Declarative YAML configuration** for cluster topology and options
- **Script and config generation** for all nodes (primary and replicas)
- **Systemd integration** for health checks and failover automation
- **Automatic failover** to replica on primary failure
- **Observability**: Datadog integration and extensible monitoring
- **Cross-platform**: Linux, macOS, and Windows support

## Quickstart

1. **Install**
  ```sh
  go install github.com/HasithDeAlwis/ha-syncgen@latest
  # or build from source
  go build -o syncgen ./
  ```
2. **Create/Edit your config**
  See [`configuration.mdx`](./docs/content/docs/configuration.mdx) for all options and examples.
3. **Validate your config**
  ```sh
  syncgen validate cluster.yaml
  ```
4. **Generate scripts**
  ```sh
  syncgen build cluster.yaml
  ```
5. **Deploy**
  - Copy scripts from `generated/` to your servers
  - Run setup scripts and enable systemd services

## CLI Reference
See [CLI Reference](./docs/content/docs/cli-reference.mdx) for all commands, flags, and usage examples.

## Documentation
- [Configuration Guide](./docs/content/docs/configuration.mdx)
- [Examples](./docs/content/docs/examples.mdx)
- [Generated Scripts](./docs/content/docs/generated-scripts.mdx)
- [Troubleshooting](./docs/content/docs/troubleshooting.mdx)
- [Roadmap](./docs/content/docs/roadmap.mdx)

## DevOps & Advanced Usage
For AWS/dev automation and advanced workflows, see [`devenv/aws/main.go`](./devenv/aws/main.go).

## Contributing
- Edit docs in `docs/content/docs/`
- PRs welcome!

## License
Apache 2.0

---

## Features

- **Declarative YAML configuration** — Define your cluster, nodes, and options in a single file
- **Script & config generation** — All necessary setup, health, and failover scripts are generated for you
- **Systemd integration** — Health checks and failover logic run as systemd services/timers
- **Automatic failover** — Detects primary failure and promotes a replica
- **Observability** — Datadog integration and extensible monitoring hooks
- **Cross-platform** — Works on Linux, macOS, and Windows

## Roadmap

- [ ] `syncgen deploy` command for automated deployment
- [ ] `syncgen status` command for cluster monitoring
- [ ] `syncgen promote` command for manual failover
- [ ] Advanced validation and error handling
- [ ] Integration tests with real PostgreSQL instances
- [ ] Docker support for development/testing
