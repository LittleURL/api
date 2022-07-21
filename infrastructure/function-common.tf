locals {
  envvar_queues = {
    "QUEUES_USERUPDATE" = aws_sqs_queue.user_update.name
  }
  envvar_tables = {
    "TABLES_DOMAINS" = aws_dynamodb_table.domains.id
  }
  envvar_default = merge(local.envvar_queues, local.envvar_tables, {
    "LUMIGO_USE_TRACER_EXTENSION" = var.lumigo_token == "" ? false : true
  })
}

# ----------------------------------------------------------------------------------------------------------------------
# Function deployment package storage
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_s3_bucket" "functions" {
  bucket_prefix = "${local.prefix}function-deployment"

  tags = {
    internal = true
  }
}

resource "aws_s3_bucket_acl" "functions" {
  bucket = aws_s3_bucket.functions.id
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "functions" {
  bucket = aws_s3_bucket.functions.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "functions" {
  bucket = aws_s3_bucket.functions.bucket

  rule {
    bucket_key_enabled = true
    apply_server_side_encryption_by_default {
      sse_algorithm = "aws:kms"
    }
  }
}
