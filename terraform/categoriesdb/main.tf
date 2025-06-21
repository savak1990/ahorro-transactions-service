locals {
  score_lsi_name = "${var.db_table_name}-score-lsi"
}

resource "aws_dynamodb_table" "categories" {
  name         = var.db_table_name
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "user_id"
  range_key = "score"

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "score"
    type = "N"
  }

  point_in_time_recovery {
    enabled = true
  }
}
