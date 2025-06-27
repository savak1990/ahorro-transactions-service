output "db_endpoint" {
  description = "The writer endpoint of the Aurora cluster."
  value       = aws_rds_cluster.aurora.endpoint
}

output "db_reader_endpoint" {
  description = "The reader endpoint of the Aurora cluster."
  value       = aws_rds_cluster.aurora.reader_endpoint
}

output "db_port" {
  description = "The port of the Aurora database."
  value       = aws_rds_cluster.aurora.port
}

output "db_name" {
  description = "The database name."
  value       = aws_rds_cluster.aurora.database_name
}
