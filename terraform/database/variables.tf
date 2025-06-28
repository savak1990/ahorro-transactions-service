variable "db_identifier" {
  description = "The identifier for the RDS instance."
  type        = string
}

variable "engine_version" {
  description = "PostgreSQL engine version."
  type        = string
  default     = "16.8"
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

variable "instance_class" {
  description = "Instance class for RDS instance. db.t3.micro is the cheapest option."
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  description = "Initial allocated storage in GB."
  type        = number
  default     = 20 # Minimum for PostgreSQL
}

variable "max_allocated_storage" {
  description = "Maximum allocated storage in GB for autoscaling."
  type        = number
  default     = 100
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
