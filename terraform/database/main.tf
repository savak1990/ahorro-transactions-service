# Regular RDS PostgreSQL instance for Ahorro Transactions Service
# Optimized for cost-effective testing

resource "aws_db_instance" "postgres" {
  identifier = var.db_identifier

  # Database engine
  engine         = "postgres"
  engine_version = var.engine_version

  # Instance configuration - using smallest available instance for cost optimization
  instance_class        = var.instance_class
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage
  storage_type          = "gp3" # General Purpose SSD (gp3) is more cost effective than gp2
  storage_encrypted     = false # Disable encryption to reduce costs for testing

  # Database configuration
  db_name  = var.db_name
  username = var.master_username
  password = var.master_password
  port     = 5432

  # Network configuration
  db_subnet_group_name   = aws_db_subnet_group.postgres_subnet_group.name
  vpc_security_group_ids = [aws_security_group.postgres_sg.id]
  publicly_accessible    = var.enable_public_access

  # Backup and maintenance
  backup_retention_period  = 1 # Minimum required (1 day) to reduce costs
  backup_window            = "02:00-04:00"
  delete_automated_backups = true
  skip_final_snapshot      = true
  deletion_protection      = false # Allow easy deletion for testing

  # Performance and monitoring
  performance_insights_enabled = false # Disable to reduce costs
  monitoring_interval          = 0     # Disable enhanced monitoring

  # Disable all database upgrades for stability
  allow_major_version_upgrade = false
  auto_minor_version_upgrade  = false

  tags = {
    Name        = var.db_identifier
    Environment = "testing"
    Purpose     = "cost-optimized"
  }
}

resource "aws_db_subnet_group" "postgres_subnet_group" {
  name       = "${var.db_identifier}-subnet-group"
  subnet_ids = var.subnet_ids

  tags = {
    Name = "${var.db_identifier}-subnet-group"
  }
}

resource "aws_security_group" "postgres_sg" {
  name        = "${var.db_identifier}-sg"
  description = "PostgreSQL RDS access for Lambdas"
  vpc_id      = var.vpc_id

  # Lambda access (always allowed)
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = var.lambda_cidr_blocks
    description = "Lambda access"
  }

  # Conditional public access for development
  dynamic "ingress" {
    for_each = var.enable_public_access ? [1] : []
    content {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      cidr_blocks = var.allowed_public_cidr_blocks
      description = "Temporary public access for development"
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name         = "${var.db_identifier}-sg"
    PublicAccess = var.enable_public_access ? "enabled" : "disabled"
  }
}
