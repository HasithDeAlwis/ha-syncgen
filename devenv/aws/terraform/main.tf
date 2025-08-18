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
  region  = "us-east-2"
  profile = var.ec2_profile
}

resource "aws_security_group" "postgres_sg" {
  name        = "postgres_sg"
  description = "Allow PostgreSQL traffic for ha-syncgen testing"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Allow all traffic (not secure for production)
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

locals {
  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    amazon-linux-extras install docker
  EOF

  instances = {
    primary = {
      role        = "primary",
      name        = "ha-syncgen-primary",
      db_user     = "admin",
      db_password = "admin_password",
      db_name     = "primary"
    }
    replica1 = {
      role        = "replica",
      name        = "ha-syncgen-replica-1",
      db_user     = "replica1_admin",
      db_password = "replica1_password",
      db_name     = "replica1"
    }
    replica2 = {
      role        = "replica",
      name        = "ha-syncgen-replica-2",
      db_user     = "replica2_admin",
      db_password = "replica2_password",
      db_name     = "replica2"
    }
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "ha-syncgen-vpc"
  }
}

resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-east-2a"
}

resource "aws_instance" "ha_syncgen" {
  for_each = local.instances

  ami                    = "ami-0169aa51f6faf20d5"
  instance_type          = "t2.micro"
  key_name               = var.key_pair_name
  subnet_id              = aws_subnet.public_subnet.id
  vpc_security_group_ids = [aws_security_group.postgres_sg.id]
  user_data              = local.user_data

  tags = {
    Name    = each.value.name
    Role    = each.value.role
    Purpose = "Testing ha-syncgen PostgreSQL HA"
  }
}

output "instance_details" {
  value = {
    for instance in sort(keys(aws_instance.ha_syncgen)) : instance => {
      ip_address  = aws_instance.ha_syncgen[instance].public_ip,
      role        = local.instances[instance].role,
      db_user     = local.instances[instance].db_user,
      db_password = local.instances[instance].db_password,
      db_name     = local.instances[instance].db_name
    }
  }
}

output "datadog_details" {
  value = {
    api_key = var.datadog_api_key,
    site    = var.datadog_site
  }
}
