locals {
  function_name_cognito_post_authentication = "${var.prefix}cognito-post-authentication"
}

# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_cogntio_post_authentication" {
  source = "../modules/lambda-function"

  aws_account = var.aws_account
  aws_region  = var.aws_region

  name          = local.function_name_cognito_post_authentication
  source_key    = "cogntio-post-authentication.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {
    "COGNITOPOOLID" = var.cognito_pool_id
  })
}

resource "aws_lambda_permission" "cogntio_post_authentication" {
  statement_id  = "AllowCognitoInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_cogntio_post_authentication.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = var.cognito_pool_arn
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role_policy" "lambda_cogntio_post_authentication_cognito" {
  name   = "Cogntio"
  role   = module.lambda_cogntio_post_authentication.role_id
  policy = data.aws_iam_policy_document.lambda_cogntio_post_authentication_cognito.json
}

data "aws_iam_policy_document" "lambda_cogntio_post_authentication_cognito" {
  statement {
    sid = "CognitoAdminUsers"

    actions = [
      "cognito-idp:ListUsers",
      "cognito-idp:AdminUpdateUserAttributes"
    ]

    resources = [var.cognito_pool_arn]
  }
}
