output "api_url" {
  value     = module.ahorro_transactions_service.service_url
  sensitive = true
}

output "cognito_user_pool_arn" {
  description = "value of the Cognito User Pool ARN"
  value       = data.terraform_remote_state.cognito.outputs.user_pool_arn
}
