locals {
  openapi_template_path = "${path.module}/schema/openapi.yml.tmpl"
  lambda_name           = "${var.base_name}-app-lambda"
}

resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.app.function_name}"
  retention_in_days = 14
}

data "aws_s3_object" "lambda_zip" {
  bucket = var.app_s3_bucket_name
  key    = var.app_s3_artifact_zip_key
}

resource "aws_lambda_function" "app" {
  function_name     = local.lambda_name
  role              = aws_iam_role.lambda_role.arn
  handler           = "bootstrap"
  runtime           = "provided.al2"
  s3_bucket         = var.app_s3_bucket_name
  s3_key            = var.app_s3_artifact_zip_key
  s3_object_version = data.aws_s3_object.lambda_zip.version_id
  source_code_hash  = data.aws_s3_object.lambda_zip.etag
  timeout           = 30

  vpc_config {
    subnet_ids         = var.lambda_subnet_ids
    security_group_ids = [aws_security_group.lambda_sg.id]
  }

  environment {
    variables = {
      # Aurora PostgreSQL connection details
      DB_HOST     = var.db_endpoint
      DB_NAME     = var.db_name
      DB_USER     = var.db_username
      DB_PASSWORD = var.db_password
      DB_PORT     = "5432"

      # Exchange Rate Source (API Key and Cache Dynamo DB)
      EXCHANGE_RATE_DB_NAME = var.exchange_rate_db_name

      # SSL Configuration
      SSL_MODE = "require"

      # Application configuration
      LOG_LEVEL = var.log_level
    }
  }
}

resource "aws_security_group" "lambda_sg" {
  name        = "${local.lambda_name}-sg"
  description = "Security group for Lambda function"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# VPC Endpoint for DynamoDB
resource "aws_vpc_endpoint" "dynamodb" {
  vpc_id            = var.vpc_id
  service_name      = "com.amazonaws.eu-west-1.dynamodb"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = data.aws_route_tables.vpc_route_tables.ids

  tags = {
    Name = "${var.base_name}-dynamodb-endpoint"
  }
}

# Get route tables for the VPC
data "aws_route_tables" "vpc_route_tables" {
  vpc_id = var.vpc_id
}

resource "aws_iam_role" "lambda_role" {
  name        = "${local.lambda_name}-role"
  description = "IAM Role for ${local.lambda_name}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Principal = {
          Service = ["lambda.amazonaws.com", "apigateway.amazonaws.com"]
        }
        Effect = "Allow"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "app_lambda_exec_policy_attachment" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_vpc_policy_attachment" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "logging_policy" {
  name = "${local.lambda_name}-logging-policy"
  role = aws_iam_role.lambda_role.name

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Action" : [
          "logs:CreateLogStream",
          "logs:DescribeLogStream",
          "logs:PutLogEvents"
        ],
        "Resource" : "arn:aws:logs:*:*:log-group:/aws/lambda/${local.lambda_name}*",
        "Effect" : "Allow"
      }
    ]
  })
}

resource "aws_iam_role_policy" "cloudwatch_policy" {
  name = "${local.lambda_name}-cloudwatch"
  role = aws_iam_role.lambda_role.name

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Action" : [
          "cloudwatch:PutMetricData"
        ],
        "Effect" : "Allow",
        "Resource" : "*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_s3_access" {
  name = "${local.lambda_name}-s3-access"
  role = aws_iam_role.lambda_role.name

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:GetObject"
        ],
        Resource = "arn:aws:s3:::ahorro-artifacts/*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_dynamodb_access" {
  name = "${local.lambda_name}-dynamodb-access"
  role = aws_iam_role.lambda_role.name

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem"
        ],
        Resource = "arn:aws:dynamodb:*:*:table/${var.exchange_rate_db_name}"
      }
    ]
  })
}

module "apigateway" {
  source                      = "github.com/savak1990/ahorro-shared//terraform/apigateway_http?ref=v0.0.1"
  api_name                    = var.api_name
  domain_name                 = var.domain_name
  zone_id                     = var.zone_id
  certificate_arn             = var.certificate_arn
  lambda_name                 = aws_lambda_function.app.function_name
  lambda_invoke_arn           = aws_lambda_function.app.invoke_arn
  cognito_user_pool_id        = var.cognito_user_pool_id
  cognito_user_pool_client_id = var.cognito_user_pool_client_id
  cognito_auth_paths = [
    "ANY /transactions",
    "ANY /transactions/{transaction_id}",
    "ANY /balances",
    "ANY /balances/{balance_id}",
    "ANY /categories",
    "ANY /categories/{category_id}",
    "ANY /merchants",
    "ANY /merchants/{merchant_id}",
  ]
  cognito_unauth_paths = [
    "GET /info",
    "GET /health",
    "GET /docs",
    "GET /schema",
    "GET /schema/info",
    "GET /schema/raw",
    "GET /schema/json",
  ]
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  required_version = ">= 1.0"
}
