# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_lambda_function" "default" {
  function_name = var.name
  role          = aws_iam_role.lambda_execution.arn

  # handler
  handler   = "bootstrap"
  s3_bucket = var.source_bucket
  s3_key    = coalesce(var.source_key, "${var.name}.zip")

  # runtime
  runtime       = "provided.al2"
  architectures = var.architectures
  timeout       = var.timeout
  memory_size   = var.memory

  environment {
    variables = var.environment_variables
  }

  depends_on = [
    aws_s3_object.placeholder
  ]
}

# ----------------------------------------------------------------------------------------------------------------------
# Permssion
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role" "lambda_execution" {
  name               = var.name
  path               = var.role_path
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
}

data "aws_iam_policy_document" "lambda_assume" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Placeholder object (race condition when deploying for the first time)
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_s3_object" "placeholder" {
  bucket = var.source_bucket
  key    = coalesce(var.source_key, "${var.name}.zip")

  content_type = "application/zip"
  content_base64 = join("", [
    "UEsDBAoAAAAAABB/91TtxOLVFgAAABYAAAAJABwAYm9vdHN0cm",
    "FwVVQJAAMvGtxiLhrcYnV4CwABBOgDAAAE6AMAAFRoaXMgaXMg",
    "YSBwbGFjZWhvbGRlci5QSwECHgMKAAAAAAAQf/dU7cTi1RYAAA",
    "AWAAAACQAYAAAAAAABAAAAtIEAAAAAYm9vdHN0cmFwVVQFAAMv",
    "GtxidXgLAAEE6AMAAAToAwAAUEsFBgAAAAABAAEATwAAAFkAAA",
    "AAAA=="
  ])

  lifecycle {
    ignore_changes = [
      etag,
      tags,
      tags_all
    ]
  }
}
