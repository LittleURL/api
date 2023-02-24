output "function_arn_cognito_post_authentication" {
  # manually defined ARN due to terraform dependency cycle
  value = "arn:aws:lambda:${var.aws_region}:${var.aws_account}:function:${local.function_name_cognito_post_authentication}"
}

output "function_arn_cognito_pre_token_gen" {
  # manually defined ARN due to terraform dependency cycle
  value = "arn:aws:lambda:${var.aws_region}:${var.aws_account}:function:${local.function_name_cognito_pre_token_gen}"
}

output "functions_bucket" {
  value = aws_s3_bucket.functions.id
}