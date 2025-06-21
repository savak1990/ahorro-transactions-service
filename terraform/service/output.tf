output "service_url" {
  description = "The URL of the ahorro transactions service."
  value       = module.apigateway.api_gateway_url
}
