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

  environment {
    variables = var.environment_variables
  }
}

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
