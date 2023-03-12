locals {
  envvar_tables = {
    "Tables_Domains"     = var.ddb_table_names.domains
    "Tables_Users"       = var.ddb_table_names.users
    "Tables_UserRoles"   = var.ddb_table_names.user_roles
    "Tables_UserInvites" = var.ddb_table_names.user_invites
    "Tables_Links"       = var.ddb_table_names.links
  }


  envvar_lumigo = var.lumigo_token == "" ? { "LUMIGO_USE_TRACER_EXTENSION" = false } : {
    "LUMIGO_USE_TRACER_EXTENSION" = true,
    "LUMIGO_TRACER_TOKEN"         = var.lumigo_token
  }

  // merge env vars and convert keys to uppercase (required by configor)
  envvar_default = { for k, v in merge(local.envvar_tables, local.envvar_lumigo, var.environment) : upper(k) => v }
}

# ----------------------------------------------------------------------------------------------------------------------
# Function deployment package storage
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_s3_bucket" "functions" {
  bucket_prefix = "${var.prefix}function-deployment-"

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

resource "aws_s3_bucket_public_access_block" "functions" {
  bucket = aws_s3_bucket.functions.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "functions" {
  bucket = aws_s3_bucket.functions.id

  rule {
    id = "delete-old-versions"

    noncurrent_version_transition {
      noncurrent_days = 3
      storage_class   = "GLACIER"
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }

    status = "Enabled"
  }

  depends_on = [aws_s3_bucket_versioning.functions]
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
      "arn:aws:lambda:${var.aws_region}:${var.aws_account}:function:${var.prefix}*"
    ]
  }
}
