# README: Cloud Testing with Real PostgreSQL Instances

This directory contains Terraform configuration and scripts to deploy **real PostgreSQL instances** in AWS and test your `ha-syncgen` tool with actual cloud infrastructure.

## 🚀 What This Does

1. **Deploys 3 EC2 instances** in AWS with PostgreSQL 14 installed
2. **Generates ha-syncgen scripts** for the real IP addresses
3. **Executes your generated scripts** on the actual servers
4. **Tests streaming replication** between real PostgreSQL instances
5. **Proves your tool works** in production cloud environments

## 📋 Prerequisites

```bash
# Install Terraform
brew install terraform

# Configure AWS CLI
aws configure

# Ensure you have AWS credentials with EC2 permissions
```

## 🏃‍♂️ Quick Start

```bash
# Navigate to terraform directory
cd terraform

# Make script executable
chmod +x deploy-and-test.sh

# Deploy and test (will prompt for confirmation)
./deploy-and-test.sh
```

## 📁 Files

- `main.tf` - Terraform configuration for AWS infrastructure
- `aws-cluster.yaml` - Template cluster configuration  
- `deploy-and-test.sh` - Complete deployment and testing script
- `README.md` - This file

## 💰 Cost Considerations

- Uses `t3.micro` instances (free tier eligible)
- Estimated cost: ~$0.50/hour for 3 instances
- **Remember to destroy resources** when done testing

## 🧪 What Gets Tested

### ✅ Real Infrastructure
- 3 EC2 instances with PostgreSQL 14
- Proper VPC, subnets, and security groups
- Real network connectivity between instances

### ✅ Generated Scripts Execution
- `setup_primary.sh` runs on primary server
- `setup_replication.sh` runs on both replicas
- `health_check.sh` can be tested manually

### ✅ PostgreSQL Streaming Replication
- Replication slots created and used
- pg_basebackup executed from replicas
- Data replication verified with test queries

### ✅ Configuration Files Applied
- `postgresql.conf.patch` applied to primary
- `pg_hba.conf.patch` applied for authentication
- Systemd service files created

## 🔍 Manual Testing Steps

After deployment, you can manually test:

```bash
# SSH to primary
ssh -i ~/.ssh/ha-syncgen-test.pem ec2-user@<PRIMARY_IP>

# Check PostgreSQL status
sudo systemctl status postgresql-14
sudo -u postgres psql -c "SELECT * FROM pg_stat_replication;"

# SSH to replica
ssh -i ~/.ssh/ha-syncgen-test.pem ec2-user@<REPLICA_IP>

# Check replication status
sudo -u postgres psql -c "SELECT pg_is_in_recovery();"
```

## 🧹 Cleanup

```bash
# Destroy AWS resources
cd terraform
terraform destroy
```

## 📊 Results

This test proves:
- ✅ Your `ha-syncgen` tool generates **working scripts**
- ✅ Scripts work on **real Linux servers** (Amazon Linux 2)
- ✅ PostgreSQL streaming replication **actually functions**
- ✅ Configuration files **properly configure** PostgreSQL
- ✅ Your tool creates **production-ready** infrastructure automation

## 🏆 Resume Value

Successfully completing this test demonstrates:
- **Infrastructure as Code** with Terraform
- **PostgreSQL High Availability** implementation
- **AWS Cloud Deployment** experience
- **Automation Tool Development** with real-world validation
- **DevOps/SRE Skills** with database replication

This is exactly the kind of project that shows you can build tools that work in production cloud environments!
