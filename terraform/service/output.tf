output "service_url" {
  description = "The URL of the ahorro transactions service."
  value       = module.apigateway.api_gateway_url
}

output "api_url" {
  description = "The URL of the transactions service API."
  value       = module.apigateway.api_gateway_url
}

output "api_gateway_id" {
  description = "The ID of the API Gateway."
  value       = module.apigateway.http_api_id
}

output "lambda_function_name" {
  description = "The name of the Lambda function."
  value       = aws_lambda_function.app.function_name
}
