variable "cluster_identifier" {
  description = "The identifier for the Aurora cluster."
  type        = string
}

variable "engine_version" {
  description = "Aurora PostgreSQL engine version."
  type        = string
  default     = "15.3"
}

variable "db_name" {
  description = "The name of the database to create."
  type        = string
}

variable "master_username" {
  description = "Master username for the database. Should be sourced from your existing Secrets Manager secret."
  type        = string
}

variable "master_password" {
  description = "Master password for the database. Should be sourced from your existing Secrets Manager secret."
  type        = string
  sensitive   = true
}

variable "min_capacity" {
  description = "Minimum Aurora Serverless v2 capacity (ACUs)."
  type        = number
  default     = 0.5
}

variable "max_capacity" {
  description = "Maximum Aurora Serverless v2 capacity (ACUs)."
  type        = number
  default     = 2
}

variable "instance_count" {
  description = "Number of RDS cluster instances."
  type        = number
  default     = 1
}

variable "instance_class" {
  description = "Instance class for Aurora cluster instances."
  type        = string
  default     = "db.serverless"
}

variable "subnet_ids" {
  description = "List of subnet IDs for the DB subnet group."
  type        = list(string)
}

variable "vpc_id" {
  description = "VPC ID for the Aurora cluster."
  type        = string
}

variable "lambda_cidr_blocks" {
  description = "List of CIDR blocks allowed to access Aurora (e.g., Lambda)."
  type        = list(string)
}

variable "enable_public_access" {
  description = "Enable public access to Aurora for development. Set to true only for temporary access."
  type        = bool
  default     = false
}

variable "allowed_public_cidr_blocks" {
  description = "CIDR blocks allowed for public access (your IP). Only used when enable_public_access is true."
  type        = list(string)
  default     = []
}
