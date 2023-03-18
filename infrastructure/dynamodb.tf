resource "aws_dynamodb_table" "domains" {
  name                        = "${local.prefix}domains"
  billing_mode                = "PAY_PER_REQUEST"
  deletion_protection_enabled = true

  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_dynamodb_table" "links" {
  name                        = "${local.prefix}links"
  billing_mode                = "PAY_PER_REQUEST"
  deletion_protection_enabled = true

  hash_key  = "domain_id"
  range_key = "uri"

  attribute {
    name = "domain_id"
    type = "S"
  }

  attribute {
    name = "uri"
    type = "S"
  }

  attribute {
    name = "updated_at"
    type = "N"
  }

  ttl {
    attribute_name = "expires_at"
    enabled        = true
  }

  local_secondary_index {
    name            = "updated"
    range_key       = "updated_at"
    projection_type = "ALL"
  }

  lifecycle {
    prevent_destroy = true
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Auth
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_dynamodb_table" "user_roles" {
  name                        = "${local.prefix}user-roles"
  billing_mode                = "PAY_PER_REQUEST"
  deletion_protection_enabled = true

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

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_dynamodb_table" "user_invites" {
  name                        = "${local.prefix}user-invites"
  billing_mode                = "PAY_PER_REQUEST"
  deletion_protection_enabled = true

  hash_key = "id"

  ttl {
    enabled        = true
    attribute_name = "expires_at"
  }

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "domain_id"
    type = "S"
  }

  global_secondary_index {
    name            = "domain-invites"
    hash_key        = "domain_id"
    range_key       = "id"
    projection_type = "ALL"
  }

  lifecycle {
    prevent_destroy = true
  }
}

# This mostly exists to make querying easier because Cognito's APIs are dogshit
resource "aws_dynamodb_table" "users" {
  name                        = "${local.prefix}users"
  billing_mode                = "PAY_PER_REQUEST"
  deletion_protection_enabled = true

  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }

  lifecycle {
    prevent_destroy = true
  }
}
