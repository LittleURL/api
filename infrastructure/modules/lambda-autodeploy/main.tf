resource "aws_lambda_function" "deployment_trigger" {
  function_name = var.name
  role          = aws_iam_role.lambda_execution.arn

  # handler
  handler          = "handler.handler"
  filename         = data.archive_file.deployment_trigger_src.output_path
  source_code_hash = data.archive_file.deployment_trigger_src.output_base64sha256

  # runtime
  runtime       = "nodejs18.x"
  architectures = ["arm64"]
  timeout       = 5
  memory_size   = 128

  environment {
    variables = {
      "FUNCTION_NAME_PREFIX" = var.function_name_prefix
    }
  }
}

data "archive_file" "deployment_trigger_src" {
  type        = "zip"
  source_file = "${path.module}/handler.js"
  output_path = "${path.module}/handler.zip"
}

# ----------------------------------------------------------------------------------------------------------------------
# Role
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role" "lambda_execution" {
  name               = var.name
  path               = var.role_path
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role_policy" "lambda" {
  name   = "Lambda"
  role   = aws_iam_role.lambda_execution.id
  policy = data.aws_iam_policy_document.lambda.json
}

data "aws_iam_policy_document" "lambda" {
  statement {
    sid = "GetFunction"

    actions = [
      "lambda:GetFunction",
      "lambda:ListFunctions"
    ]

    resources = ["*"]
  }

  statement {
    sid = "UpdateFunctions"

    actions = [
      "lambda:UpdateFunctionCode"
    ]

    resources = [
      "*"
    ]

    condition {
      test     = "StringEquals"
      variable = "aws:ResourceTag/littleurl-autodeploy-enabled"
      values   = ["true"]
    }
  }

  statement {
    sid = "S3"

    actions = [
      "s3:GetObject"
    ]

    resources = [
      "arn:aws:s3:::${var.bucket}/*"
    ]
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Trigger
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.deployment_trigger.arn
  principal     = "s3.amazonaws.com"
  source_arn    = "arn:aws:s3:::${var.bucket}"
}

resource "aws_s3_bucket_notification" "deployment_trigger" {
  bucket = var.bucket

  lambda_function {
    lambda_function_arn = aws_lambda_function.deployment_trigger.arn
    events              = ["s3:ObjectCreated:*"]
    filter_suffix       = ".zip"
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}
