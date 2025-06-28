output "api_url" {
  value     = module.ahorro_transactions_service.service_url
  sensitive = true
}

output "cognito_user_pool_arn" {
  description = "value of the Cognito User Pool ARN"
  value       = data.terraform_remote_state.cognito.outputs.user_pool_arn
}

output "db_name" {
  description = "The name of the PostgreSQL database."
  value       = module.transactions_db.db_name
}

output "db_endpoint" {
  description = "The endpoint of the PostgreSQL database."
  value       = module.transactions_db.db_endpoint
}

output "db_port" {
  description = "The port of the PostgreSQL database."
  value       = module.transactions_db.db_port
}

output "db_identifier" {
  description = "The RDS instance identifier."
  value       = module.transactions_db.db_identifier
}
