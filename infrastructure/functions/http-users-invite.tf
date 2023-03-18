# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_users_invite" {
  source = "../modules/lambda-function"

  aws_account       = var.aws_account
  aws_region        = var.aws_region
  enable_autodeploy = var.enable_autodeploy

  name          = "${var.prefix}http-users-invite"
  source_key    = "http-users-invite.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_users_invite" {
  source = "../modules/lambda-gateway"
  method = "POST"
  path   = "/domains/{domainId}/users/invite"

  function_name       = module.lambda_http_users_invite.function_name
  function_invoke_arn = module.lambda_http_users_invite.function_invoke_arn

  gateway_id            = var.gateway_id
  gateway_execution_arn = var.gateway_execution_arn

  authorizer_id = var.authorizer_id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_users_invite_dynamodb" {
  source = "../modules/iam-dynamodb"
  role   = module.lambda_http_users_invite.role_id

  tables = [
    {
      arn          = var.ddb_table_arns.user_invites
      enable_write = true
    },
    {
      arn          = var.ddb_table_arns.user_roles
      enable_read = true
    },
    {
      arn          = var.ddb_table_arns.domains
      enable_read = true
    }
  ]
}

resource "aws_iam_role_policy" "lambda_http_users_invite_ses" {
  name   = "SES"
  role   = module.lambda_http_users_invite.role_id
  policy = data.aws_iam_policy_document.lambda_http_users_invite_ses.json
}
data "aws_iam_policy_document" "lambda_http_users_invite_ses" {
  statement {
    sid = "SendEmail"

    actions = [
      "ses:SendEmail",
      "ses:SendRawEmail"
    ]

    resources = ["*"]

    condition {
      test     = "StringLike"
      variable = "ses:FromAddress"
      values   = var.email_allowed_from_addresses
    }
  }
}
