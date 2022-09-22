resource "aws_dynamodb_table" "domains" {
  name         = "${local.prefix}domains"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "domain_users" {
  name         = "${local.prefix}domain-users"
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "domain_id"
  range_key = "user_id"

  attribute {
    name = "domain_id"
    type = "S"
  }

  attribute {
    name = "user_id"
    type = "S"
  }

  global_secondary_index {
    name            = "user-domains"
    hash_key        = "user_id"
    range_key       = "domain_id"
    projection_type = "ALL"
  }
}

