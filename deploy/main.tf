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

data "terraform_remote_state" "db" {
  backend = "s3"
  config = {
    bucket         = "ahorro-app-state"
    key            = "stable/transactions-service/db/terraform.tfstate"
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

locals {
  base_name               = "${var.app_name}-${var.service_name}-${var.env}"
  app_lambda_name         = "${local.base_name}-app-lambda"
  app_s3_bucket_name      = "ahorro-artifacts"
  app_s3_artifact_zip_key = "transactions/${var.env}/transactions-lambda.zip"
  full_api_name           = "api-${local.base_name}"

  db_subnet_ids = data.aws_subnets.default.ids
  vpc_id        = data.aws_vpc.default.id

  secret_name              = "${var.app_name}-app-secrets"
  domain_name              = jsondecode(data.aws_secretsmanager_secret_version.ahorro_app.secret_string)["domain_name"]
  transactions_db_username = jsondecode(data.aws_secretsmanager_secret_version.ahorro_app.secret_string)["transactions_db_username"]
  transactions_db_password = jsondecode(data.aws_secretsmanager_secret_version.ahorro_app.secret_string)["transactions_db_password"]
}

module "ahorro_transactions_service" {
  source = "../terraform/service"

  # VPC Configuration
  vpc_id            = local.vpc_id
  lambda_subnet_ids = local.db_subnet_ids

  # Aurora Database Configuration
  db_endpoint = data.terraform_remote_state.db.outputs.db_endpoint
  db_name     = data.terraform_remote_state.db.outputs.db_name
  db_username = local.transactions_db_username
  db_password = local.transactions_db_password

  # Application Configuration
  log_level = "info"

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
