provider "aws" {
  region = "eu-west-1"

  default_tags {
    tags = {
      Environment = "dev"
      Project     = "ahorro-app"
      Service     = "ahorro-transactions-service"
      Terraform   = "true"
    }
  }
}

data "aws_secretsmanager_secret" "ahorro_app" {
  name = local.secret_name
}

data "aws_secretsmanager_secret_version" "ahorro_app" {
  secret_id = data.aws_secretsmanager_secret.ahorro_app.id
}

data "aws_acm_certificate" "cert" {
  domain      = "*.${local.domain_name}"
  statuses    = ["ISSUED"]
  most_recent = true
}

data "aws_route53_zone" "public" {
  name         = local.domain_name
  private_zone = false
}

data "terraform_remote_state" "cognito" {
  backend = "s3"
  config = {
    bucket         = "ahorro-app-state"
    key            = "stable/cognito/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "ahorro-app-state-lock"
    encrypt        = true
  }
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Get CIDR blocks of the subnets where Lambda will run
data "aws_subnet" "lambda_subnets" {
  for_each = toset(data.aws_subnets.default.ids)
  id       = each.value
}

locals {
  base_name               = "${var.app_name}-${var.service_name}-${var.env}"
  db_name                 = "${var.app_name}_${var.service_name}_${var.env}_db"
  db_cluster_name         = "${local.db_name}-cluster"
  app_lambda_name         = "${local.base_name}-app-lambda"
  app_s3_bucket_name      = "ahorro-artifacts"
  app_s3_artifact_zip_key = "transactions/transactions-lambda.zip"
  full_api_name           = "api-${local.base_name}"

  db_subnet_ids = data.aws_subnets.default.ids
  vpc_id        = data.aws_vpc.default.id
  # Extract CIDR blocks from Lambda subnets
  lambda_cidr_blocks = [for subnet in data.aws_subnet.lambda_subnets : subnet.cidr_block]

  secret_name              = "${var.app_name}-app-secrets"
  domain_name              = jsondecode(data.aws_secretsmanager_secret_version.ahorro_app.secret_string)["domain_name"]
  transactions_db_username = jsondecode(data.aws_secretsmanager_secret_version.ahorro_app.secret_string)["transactions_db_username"]
  transactions_db_password = jsondecode(data.aws_secretsmanager_secret_version.ahorro_app.secret_string)["transactions_db_password"]
}

module "transactions_db" {
  source = "../terraform/database"

  db_identifier   = "${local.base_name}-db"
  db_name         = local.db_name
  engine_version  = "16.8" # Latest stable PostgreSQL
  master_username = local.transactions_db_username
  master_password = local.transactions_db_password

  # Cost-optimized settings
  instance_class        = "db.t3.micro" # Cheapest option
  allocated_storage     = 20            # Minimum storage
  max_allocated_storage = 50            # Small autoscaling limit

  # Network configuration
  subnet_ids         = local.db_subnet_ids
  vpc_id             = local.vpc_id
  lambda_cidr_blocks = local.lambda_cidr_blocks

  # Temporary public access configuration
  enable_public_access       = var.enable_db_public_access
  allowed_public_cidr_blocks = var.enable_db_public_access ? [var.my_ip_cidr] : []
}

module "ahorro_transactions_service" {
  source = "../terraform/service"

  # VPC Configuration
  vpc_id            = local.vpc_id
  lambda_subnet_ids = local.db_subnet_ids

  # Aurora Database Configuration
  db_endpoint = module.transactions_db.db_endpoint
  db_name     = module.transactions_db.db_name
  db_username = local.transactions_db_username
  db_password = local.transactions_db_password

  # API Gateway Configuration
  api_name                    = local.full_api_name
  domain_name                 = local.domain_name
  base_name                   = local.base_name
  app_s3_bucket_name          = local.app_s3_bucket_name
  app_s3_artifact_zip_key     = local.app_s3_artifact_zip_key
  certificate_arn             = data.aws_acm_certificate.cert.arn
  zone_id                     = data.aws_route53_zone.public.zone_id
  cognito_user_pool_id        = data.terraform_remote_state.cognito.outputs.user_pool_id
  cognito_user_pool_client_id = data.terraform_remote_state.cognito.outputs.user_pool_client_id

  depends_on = [module.transactions_db]
}

terraform {
  backend "s3" {
    bucket = "ahorro-app-state"
    ### Please update "savak" to your user name if you're going to try deploying this yourself
    key            = "dev/ahorro-transactions-service/savak/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "ahorro-app-state-lock"
    encrypt        = true
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
