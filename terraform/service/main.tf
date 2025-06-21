locals {
  openapi_template_path = "${path.module}/schema/openapi.yml.tmpl"
  lambda_name           = "${var.base_name}-app-lambda"
}

resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.app.function_name}"
  retention_in_days = 14
}

resource "aws_lambda_function" "app" {
  function_name    = local.lambda_name
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2"
  filename         = var.app_handler_zip
  source_code_hash = filebase64sha256(var.app_handler_zip)

  environment {
    variables = {
      TRANSACTIONS_DYNAMODB_TABLE = var.transactions_db_table_name
      CATEGORIES_DYNAMODB_TABLE   = var.categories_db_table_name
    }
  }
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

resource "aws_iam_role_policy" "categories_dynamodb_policy" {
  name = "categories-dynamodb-policy"
  role = aws_iam_role.lambda_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:DescribeTable"
        ],
        Effect = "Allow",
        Resource = [
          "arn:aws:dynamodb:*:*:table/${var.categories_db_table_name}",
          "arn:aws:dynamodb:*:*:table/${var.categories_db_table_name}/index/*"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy" "transactions_dynamodb_policy" {
  name = "transactions-dynamodb-policy"
  role = aws_iam_role.lambda_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:DescribeTable"
        ],
        Effect = "Allow",
        Resource = [
          "arn:aws:dynamodb:*:*:table/${var.transactions_db_table_name}",
          "arn:aws:dynamodb:*:*:table/${var.transactions_db_table_name}/index/*"
        ]
      }
    ]
  })
}

module "apigateway" {
  source                      = "../../../ahorro-shared/terraform/apigateway_http"
  api_name                    = var.api_name
  domain_name                 = var.domain_name
  zone_id                     = var.zone_id
  certificate_arn             = var.certificate_arn
  lambda_name                 = aws_lambda_function.app.function_name
  lambda_invoke_arn           = aws_lambda_function.app.invoke_arn
  cognito_user_pool_id        = var.cognito_user_pool_id
  cognito_user_pool_client_id = var.cognito_user_pool_client_id
  cognito_auth_paths = [
    "GET /transactions",
    "POST /transactions",
    "GET /transactions/{transaction_id}",
    "PUT /transactions/{transaction_id}",
    "DELETE /transactions/{transaction_id}",
    "GET /categories"
  ]
  cognito_unauth_paths = [
    "GET /info",
    "GET /health",
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
