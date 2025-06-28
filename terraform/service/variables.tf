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

variable "vpc_id" {
  description = "VPC ID for Lambda networking."
  type        = string
}

variable "lambda_subnet_ids" {
  description = "List of subnet IDs for Lambda networking."
  type        = list(string)
}

variable "db_endpoint" {
  description = "Aurora cluster endpoint."
  type        = string
}

variable "db_port" {
  description = "Aurora database port."
  type        = number
  default     = 5432
}

variable "db_name" {
  description = "Aurora database name."
  type        = string
}

variable "db_username" {
  description = "Aurora database username."
  type        = string
}

variable "db_password" {
  description = "Aurora database password."
  type        = string
  sensitive   = true
}

variable "log_level" {
  description = "Log level for the application (debug, info, warn, error)"
  type        = string
  default     = "info"
  validation {
    condition     = contains(["debug", "info", "warn", "error"], var.log_level)
    error_message = "Log level must be one of: debug, info, warn, error"
  }
}
