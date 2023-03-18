# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_invites_accept" {
  source = "../modules/lambda-function"

  aws_account       = var.aws_account
  aws_region        = var.aws_region
  enable_autodeploy = var.enable_autodeploy

  name          = "${var.prefix}http-invites-accept"
  source_key    = "http-invites-accept.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {
    "COGNITOPOOLID" = var.cognito_pool_id
  })
}

module "gateway_lambda_http_invites_accept" {
  source = "../modules/lambda-gateway"
  method = "GET"
  path   = "/invites/{inviteId}"

  function_name       = module.lambda_http_invites_accept.function_name
  function_invoke_arn = module.lambda_http_invites_accept.function_invoke_arn

  gateway_id            = var.gateway_id
  gateway_execution_arn = var.gateway_execution_arn

  authorizer_id = var.authorizer_id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_invites_accept_dynamodb" {
  source = "../modules/iam-dynamodb"
  role   = module.lambda_http_invites_accept.role_id

  tables = [
    {
      arn           = var.ddb_table_arns.user_invites
      enable_read   = true
      enable_delete = true
    },
    {
      arn          = var.ddb_table_arns.user_roles
      enable_write = true
    }
  ]
}

resource "aws_iam_role_policy" "lambda_http_invites_accept_cognito" {
  name   = "Cogntio"
  role   = module.lambda_http_invites_accept.role_id
  policy = data.aws_iam_policy_document.lambda_http_invites_accept_cognito.json
}
data "aws_iam_policy_document" "lambda_http_invites_accept_cognito" {
  statement {
    sid = "ListUsers"

    actions = [
      "cognito-idp:ListUsers"
    ]

    resources = [var.cognito_pool_arn]
  }
}
