resource "aws_dynamodb_table" "transactions" {
  name         = var.db_table_name
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "user_id"
  range_key = "transaction_id"

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "transaction_id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }
}
