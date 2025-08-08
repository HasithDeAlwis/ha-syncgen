# PostgreSQL HA Cloud Testing with Terraform
# This will create real PostgreSQL instances in AWS to test ha-syncgen

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-west-2"
}

variable "key_pair_name" {
  description = "AWS Key Pair name for EC2 access"
  type        = string
  default     = "ha-syncgen-test"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"  # Free tier eligible
}

# VPC and Networking
resource "aws_vpc" "ha_postgres_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "ha-syncgen-test-vpc"
    Purpose = "Testing PostgreSQL HA with ha-syncgen"
  }
}

resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.ha_postgres_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "${var.aws_region}a"
  map_public_ip_on_launch = true

  tags = {
    Name = "ha-syncgen-public-subnet"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.ha_postgres_vpc.id

  tags = {
    Name = "ha-syncgen-igw"
  }
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.ha_postgres_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "ha-syncgen-public-rt"
  }
}

resource "aws_route_table_association" "public_rta" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_rt.id
}

# Security Group
resource "aws_security_group" "postgres_sg" {
  name_prefix = "ha-syncgen-postgres-"
  vpc_id      = aws_vpc.ha_postgres_vpc.id

  # SSH access
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # PostgreSQL ports
  ingress {
    from_port   = 5432
    to_port     = 5435
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]  # Only within VPC
  }

  # Allow your local IP to connect to PostgreSQL (replace with your IP)
  ingress {
    from_port   = 5432
    to_port     = 5435
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # Change this to your IP for security
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "ha-syncgen-postgres-sg"
  }
}

# User data script to install PostgreSQL
locals {
  user_data = base64encode(<<-EOF
#!/bin/bash
set -e

# Update system
yum update -y

# Install PostgreSQL 14
yum install -y postgresql14-server postgresql14 postgresql14-contrib

# Install other tools
yum install -y git nc

# Initialize PostgreSQL
/usr/pgsql-14/bin/postgresql-14-setup initdb

# Enable and start PostgreSQL
systemctl enable postgresql-14
systemctl start postgresql-14

# Create postgres user password
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres123';"

# Create log directory for ha-syncgen
mkdir -p /var/log/ha-syncgen
chown postgres:postgres /var/log/ha-syncgen

# Create archive directory
mkdir -p /var/lib/postgresql/archive
chown postgres:postgres /var/lib/postgresql/archive

# Install additional tools
yum install -y wget curl

echo "PostgreSQL setup completed on $(hostname)" > /var/log/setup.log
EOF
  )
}

# EC2 Instances
resource "aws_instance" "postgres_primary" {
  ami                    = "ami-0c02fb55956c7d316"  # Amazon Linux 2 (update as needed)
  instance_type          = var.instance_type
  key_name              = var.key_pair_name
  subnet_id             = aws_subnet.public_subnet.id
  vpc_security_group_ids = [aws_security_group.postgres_sg.id]
  user_data             = local.user_data

  tags = {
    Name = "ha-syncgen-primary"
    Role = "primary"
    Purpose = "Testing ha-syncgen PostgreSQL HA"
  }
}

resource "aws_instance" "postgres_replica_1" {
  ami                    = "ami-0c02fb55956c7d316"  # Amazon Linux 2
  instance_type          = var.instance_type
  key_name              = var.key_pair_name
  subnet_id             = aws_subnet.public_subnet.id
  vpc_security_group_ids = [aws_security_group.postgres_sg.id]
  user_data             = local.user_data

  tags = {
    Name = "ha-syncgen-replica-1"
    Role = "replica"
    Purpose = "Testing ha-syncgen PostgreSQL HA"
  }
}

resource "aws_instance" "postgres_replica_2" {
  ami                    = "ami-0c02fb55956c7d316"  # Amazon Linux 2
  instance_type          = var.instance_type
  key_name              = var.key_pair_name
  subnet_id             = aws_subnet.public_subnet.id
  vpc_security_group_ids = [aws_security_group.postgres_sg.id]
  user_data             = local.user_data

  tags = {
    Name = "ha-syncgen-replica-2"
    Role = "replica"
    Purpose = "Testing ha-syncgen PostgreSQL HA"
  }
}

# Outputs
output "primary_public_ip" {
  value = aws_instance.postgres_primary.public_ip
  description = "Public IP of the primary PostgreSQL server"
}

output "primary_private_ip" {
  value = aws_instance.postgres_primary.private_ip
  description = "Private IP of the primary PostgreSQL server"
}

output "replica_1_public_ip" {
  value = aws_instance.postgres_replica_1.public_ip
  description = "Public IP of replica 1"
}

output "replica_1_private_ip" {
  value = aws_instance.postgres_replica_1.private_ip
  description = "Private IP of replica 1"
}

output "replica_2_public_ip" {
  value = aws_instance.postgres_replica_2.public_ip
  description = "Public IP of replica 2"
}

output "replica_2_private_ip" {
  value = aws_instance.postgres_replica_2.private_ip
  description = "Private IP of replica 2"
}

output "ssh_commands" {
  value = {
    primary = "ssh -i ~/.ssh/${var.key_pair_name}.pem ec2-user@${aws_instance.postgres_primary.public_ip}"
    replica_1 = "ssh -i ~/.ssh/${var.key_pair_name}.pem ec2-user@${aws_instance.postgres_replica_1.public_ip}"
    replica_2 = "ssh -i ~/.ssh/${var.key_pair_name}.pem ec2-user@${aws_instance.postgres_replica_2.public_ip}"
  }
  description = "SSH commands to connect to each instance"
}
