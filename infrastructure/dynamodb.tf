resource "aws_dynamodb_table" "domains" {
  name         = "${local.prefix}domains"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }

  stream_enabled = true
}

