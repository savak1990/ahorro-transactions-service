output "table_arn" {
  description = "The ARN of the DynamoDB table for categories."
  value       = aws_dynamodb_table.categories.arn
}
