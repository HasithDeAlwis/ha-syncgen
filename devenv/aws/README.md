
# ha-syncgen: Cloud-Ready PostgreSQL HA Automation

Automate the deployment, configuration, and testing of a real PostgreSQL High Availability (HA) cluster on AWS EC2 using Terraform, Go-based generators, and a Makefile-driven workflow.

---

## 🏗️ Architecture & Workflow Overview

- **Terraform** provisions 3 EC2 instances (1 primary, 2 replicas) in AWS.
- **Go generator** creates all Docker Compose, SQL, and deployment scripts in a single pass.
- **Makefile** orchestrates the entire workflow: infra, generation, deployment, and cleanup.

**Workflow Steps:**
1. Provision AWS infrastructure (`make aws`)
2. Generate deployment files (`make scripts`)
3. Deploy to servers (`make deploy`)
4. Generate HA sync scripts (`make syncgen`)
5. Cleanup (`make clean`)

For a detailed step-by-step workflow and implementation notes, see [WORKFLOW_SUMMARY.md](./WORKFLOW_SUMMARY.md).

---

## ⚡ Quick Start

### Prerequisites

- [Terraform](https://www.terraform.io/) installed (`brew install terraform`)
- [AWS CLI](https://aws.amazon.com/cli/) configured (`aws configure`)
- AWS credentials with EC2 permissions

### Common Commands

| Step | Command | Description |
|------|---------|-------------|
| 1 | `make init-env` | Prepare Terraform config |
| 2 | `make aws` | Deploy AWS infra (EC2, VPC, etc) |
| 3 | `make scripts` | Generate all deployment files |
| 4 | `make deploy` | Deploy Docker/PG to EC2 |
| 5 | `make syncgen` | Generate HA sync scripts |
| 6 | `make clean` | Destroy infra & clean up |

---

## �️ Makefile Targets

- `make init-env` – Prepare Terraform config
- `make aws` – Deploy AWS infra
- `make scripts` – Generate all deployment files
- `make deploy` – Deploy to EC2
- `make syncgen` – Generate HA scripts
- `make full-deploy` – Full infra + deploy
- `make full-stack` – Full infra + deploy + HA scripts
- `make dev-cycle` – Quick redeploy (scripts + deploy)
- `make clean` – Destroy infra and clean up

---

## 📁 Generated Files Structure

```
generated/
├── config.yaml
├── primary/
│   ├── docker-compose.yml
│   └── init-scripts/
│       └── 01-setup-primary.sql
├── replica1/
│   ├── docker-compose.yml
│   └── init-scripts/
│       └── 01-setup-replica1.sql
├── replica2/
│   ├── docker-compose.yml
│   └── init-scripts/
│       └── 01-setup-replica2.sql
├── deploy-to-servers.sh
└── DEPLOYMENT_README.md
```

---

## 🧪 Testing & Validation

- **Manual SSH**: Connect to EC2, check PostgreSQL status, replication
- **Automated**: All scripts run via Makefile and deployment scripts

---

## 💰 Cost Considerations

- Uses `t3.micro` instances (free tier eligible)
- Estimated cost: ~$0.50/hour for 3 instances
- **Remember to destroy resources** when done testing (`make clean`)

---

## � Troubleshooting & FAQ

- **SSH Issues**: Check your key and security group rules
- **Terraform Errors**: Ensure AWS credentials and region are set
- **Docker Issues**: Check Docker Compose logs on EC2
