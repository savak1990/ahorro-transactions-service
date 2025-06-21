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

locals {
  base_name              = "${var.app_name}-${var.service_name}-${var.env}"
  app_lambda_name        = "${local.base_name}-app-lambda"
  categories_table_name  = "${local.base_name}-categories-db"
  transaction_table_name = "${local.base_name}-transactions-db"
  secret_name            = "${var.app_name}-app-secrets"
  domain_name            = jsondecode(data.aws_secretsmanager_secret_version.ahorro_app.secret_string)["domain_name"]
  full_api_name          = "api-${local.base_name}"
}

module "categories_database" {
  source        = "../terraform/categoriesdb"
  db_table_name = local.categories_table_name
}

module "transactions_database" {
  source        = "../terraform/transactionsdb"
  db_table_name = local.transaction_table_name
}

module "ahorro_transactions_service" {
  source = "../terraform/service"

  categories_db_table_name    = local.categories_table_name
  transactions_db_table_name  = local.transaction_table_name
  api_name                    = local.full_api_name
  domain_name                 = local.domain_name
  base_name                   = local.base_name
  app_handler_zip             = var.app_handler_zip
  certificate_arn             = data.aws_acm_certificate.cert.arn
  zone_id                     = data.aws_route53_zone.public.zone_id
  cognito_user_pool_id        = data.terraform_remote_state.cognito.outputs.user_pool_id
  cognito_user_pool_client_id = data.terraform_remote_state.cognito.outputs.user_pool_client_id

  depends_on = [module.categories_database, module.transactions_database]
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
