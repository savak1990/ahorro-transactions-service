# Aurora Serverless v2 (PostgreSQL) Terraform module for Ahorro Transactions Service

resource "aws_rds_cluster" "aurora" {
  cluster_identifier = var.cluster_identifier
  engine             = "aurora-postgresql"
  engine_mode        = "provisioned"
  engine_version     = var.engine_version
  database_name      = var.db_name
  master_username    = var.master_username
  master_password    = var.master_password

  serverlessv2_scaling_configuration {
    min_capacity = var.min_capacity
    max_capacity = var.max_capacity
  }

  db_subnet_group_name        = aws_db_subnet_group.aurora_subnet_group.name
  vpc_security_group_ids      = [aws_security_group.aurora_sg.id]
  skip_final_snapshot         = true
  allow_major_version_upgrade = true # Required for major version upgrades (15.x -> 16.x)

  # Optionally, you can use the secret_arn for password rotation
  # master_user_secret { secret_arn = var.db_secret_arn }
}

resource "aws_rds_cluster_instance" "aurora_instance" {
  count               = var.instance_count
  identifier          = "${var.cluster_identifier}-instance-${count.index + 1}"
  cluster_identifier  = aws_rds_cluster.aurora.id
  instance_class      = var.instance_class
  engine              = aws_rds_cluster.aurora.engine
  engine_version      = aws_rds_cluster.aurora.engine_version
  publicly_accessible = var.enable_public_access
}

resource "aws_db_subnet_group" "aurora_subnet_group" {
  name       = "${var.cluster_identifier}-subnet-group"
  subnet_ids = var.subnet_ids
}

resource "aws_security_group" "aurora_sg" {
  name        = "${var.cluster_identifier}-sg"
  description = "Aurora access for Lambdas"
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
    Name         = "${var.cluster_identifier}-sg"
    PublicAccess = var.enable_public_access ? "enabled" : "disabled"
  }
}
