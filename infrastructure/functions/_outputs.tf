locals {
  lambda_arn_prefix = "arn:aws:lambda:${var.aws_region}:${var.aws_account}:function:"
}

output "cognito_function_arns" {
  value = {
    # manually defined ARN due to terraform dependency cycle
    post_authentication = "${local.lambda_arn_prefix}${local.function_name_cognito_post_authentication}"
    pre_token_gen       = "${local.lambda_arn_prefix}${local.function_name_cognito_pre_token_gen}"
    custom_message      = "${local.lambda_arn_prefix}${local.function_name_cognito_custom_message}"
  }
}

output "functions_bucket" {
  value = aws_s3_bucket.functions.id
}
