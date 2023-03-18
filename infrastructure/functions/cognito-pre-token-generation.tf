locals {
  function_name_cognito_pre_token_gen = "${var.prefix}cognito-pre-token-generation"
}

# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_cognito_pre_token_generation" {
  source = "../modules/lambda-function"

  aws_account       = var.aws_account
  aws_region        = var.aws_region
  enable_autodeploy = var.enable_autodeploy

  name          = local.function_name_cognito_pre_token_gen
  source_key    = "cognito-pre-token-generation.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {
    "COGNITOPOOLID" = var.cognito_pool_id
  })
}

resource "aws_lambda_permission" "cognito_pre_token_generation" {
  statement_id  = "AllowCognitoInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_cognito_pre_token_generation.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = var.cognito_pool_arn
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role_policy" "lambda_cognito_pre_token_generation_cognito" {
  name   = "Cogntio"
  role   = module.lambda_cognito_pre_token_generation.role_id
  policy = data.aws_iam_policy_document.lambda_cognito_pre_token_generation_cognito.json
}

data "aws_iam_policy_document" "lambda_cognito_pre_token_generation_cognito" {
  statement {
    sid = "CognitoAdminRead"

    actions = [
      "cognito-idp:ListUsers"
    ]

    resources = [var.cognito_pool_arn]
  }
}

module "lambda_cognito_pre_token_generation_dynamodb" {
  source = "../modules/iam-dynamodb"
  role   = module.lambda_cognito_pre_token_generation.role_id

  tables = [{
    arn          = var.ddb_table_arns.users
    enable_write = true
  }]
}
