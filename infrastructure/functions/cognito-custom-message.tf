locals {
  function_name_cognito_custom_message = "${var.prefix}cognito-custom-message"
}

# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_cognito_custom_message" {
  source = "../modules/lambda-function"

  aws_account       = var.aws_account
  aws_region        = var.aws_region
  enable_autodeploy = var.enable_autodeploy

  name          = local.function_name_cognito_custom_message
  source_key    = "cognito-custom-message.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

resource "aws_lambda_permission" "cognito_custom_message" {
  statement_id  = "AllowCognitoInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_cognito_custom_message.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = var.cognito_pool_arn
}
