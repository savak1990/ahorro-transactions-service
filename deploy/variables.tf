variable "app_name" {
  description = "The name of the application"
  type        = string
}

variable "service_name" {
  description = "The name of the service"
  type        = string
}

variable "env" {
  description = "The environment for the deployment (e.g., dev, prod, named env)"
  type        = string
}

# Temporary database public access variables
variable "enable_db_public_access" {
  description = "Enable public access to Aurora database for development. SECURITY WARNING: Only use for temporary debugging!"
  type        = bool
  default     = false
}

variable "my_ip_cidr" {
  description = "Your IP address in CIDR format (e.g., '203.0.113.1/32') for database access. Only used when enable_db_public_access is true."
  type        = string
  default     = ""
}
