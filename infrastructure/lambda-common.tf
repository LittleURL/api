locals {
  envvar_tables = {
    "TABLES_DOMAINS"   = aws_dynamodb_table.domains.id
    "TABLES_USERROLES" = aws_dynamodb_table.user_roles.id
    "TABLES_USERS"     = aws_dynamodb_table.users.id
  }
  envvar_default = merge(local.envvar_tables, local.envvar_lumigo)
}

# ----------------------------------------------------------------------------------------------------------------------
# Function deployment package storage
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_s3_bucket" "functions" {
  bucket_prefix = "${local.prefix}function-deployment-"

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

# ----------------------------------------------------------------------------------------------------------------------
# CI Permisson
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_user_policy" "function_upload" {
  name   = "DeployLambdaFunctions"
  user   = "deploy-api"
  policy = data.aws_iam_policy_document.function_upload.json
}

data "aws_iam_policy_document" "function_upload" {
  statement {
    sid = "UploadToS3"
    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:AbortMultipartUpload",
      "s3:ListBucket",
      "s3:GetObjectVersion",
      "s3:ListMultipartUploadParts"
    ]
    resources = [
      aws_s3_bucket.functions.arn,
      "${aws_s3_bucket.functions.arn}/*"
    ]
  }

  statement {
    sid = "UpdateLambdaFunction"
    actions = [
      "lambda:UpdateFunctionCode"
    ]
    resources = [
      "arn:aws:lambda:${var.aws_region}:${local.aws_account}:function:${local.prefix}*"
    ]
  }
}
