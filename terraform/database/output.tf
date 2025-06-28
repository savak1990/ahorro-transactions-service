output "db_endpoint" {
  description = "The endpoint of the RDS instance (hostname only)."
  value       = aws_db_instance.postgres.address
}

output "db_port" {
  description = "The port of the PostgreSQL database."
  value       = aws_db_instance.postgres.port
}

output "db_name" {
  description = "The database name."
  value       = aws_db_instance.postgres.db_name
}

output "db_identifier" {
  description = "The RDS instance identifier."
  value       = aws_db_instance.postgres.identifier
}
