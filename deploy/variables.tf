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
