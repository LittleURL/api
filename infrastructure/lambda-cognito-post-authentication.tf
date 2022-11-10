locals {
  function_name_cognito_post_authentication = "${local.prefix}cognito-post-authentication"
}

# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_cogntio_post_authentication" {
  source = "./modules/lambda-function"

  aws_account = local.aws_account
  aws_region  = var.aws_region

  name          = local.function_name_cognito_post_authentication
  source_key    = "cogntio-post-authentication.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {
    "COGNITOPOOLID" = aws_cognito_user_pool.main.id
  })
}

resource "aws_lambda_permission" "cogntio_post_authentication" {
  statement_id  = "AllowCognitoInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_cogntio_post_authentication.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.main.arn
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

    resources = [aws_cognito_user_pool.main.arn]
  }
}
