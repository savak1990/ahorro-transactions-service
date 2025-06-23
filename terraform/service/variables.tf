variable "transactions_db_table_name" {
  description = "The name of the DynamoDB table for transactions"
  type        = string
}

variable "categories_db_table_name" {
  description = "The name of the DynamoDB table"
  type        = string
}

variable "base_name" {
  description = "The base name for the application, used for naming resources"
  type        = string
}

variable "app_s3_bucket_name" {
  description = "The name of the S3 bucket for storing application assets"
  type        = string
}

variable "app_s3_artifact_zip_key" {
  description = "The key for the zip file in the S3 bucket"
  type        = string
}

variable "api_name" {
  description = "The name of the api. E.g. my-app-transactions-api"
  type        = string
}

variable "domain_name" {
  description = "The domain name of the API. E.g. bestapps.com"
  type        = string
}

variable "zone_id" {
  description = "The Route53 hosted zone ID for the domain."
  type        = string
}

variable "certificate_arn" {
  description = "The ARN of the ACM certificate for the custom domain."
  type        = string
}

variable "cognito_user_pool_id" {
  description = "The ID of the Cognito User Pool for authentication"
  type        = string
}

variable "cognito_user_pool_client_id" {
  description = "The ID of the Cognito User Pool Client for authentication"
  type        = string
}
