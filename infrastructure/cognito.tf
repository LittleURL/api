locals {
  cognito_domain = "auth.${local.domain}"
}

# ----------------------------------------------------------------------------------------------------------------------
# User Pool
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_cognito_user_pool" "main" {
  name = var.application

  mfa_configuration = "OPTIONAL"

  username_attributes = [
    "email"
  ]

  auto_verified_attributes = [
    "email"
  ]

  software_token_mfa_configuration {
    enabled = true
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  email_configuration {
    email_sending_account = "DEVELOPER"
    from_email_address    = local.email_from
    source_arn            = aws_ses_email_identity.noreply.arn
  }

  lambda_config {
    # manually defined ARN due to terraform dependency cycle
    pre_token_generation = "arn:aws:lambda:${var.aws_region}:${local.aws_account}:function:${local.function_name_cognito_pre_token_gen}"
    post_authentication  = "arn:aws:lambda:${var.aws_region}:${local.aws_account}:function:${local.function_name_cognito_post_authentication}"
  }
}

resource "aws_cognito_user_pool_client" "dashboard" {
  name            = "dashboard"
  user_pool_id    = aws_cognito_user_pool.main.id
  generate_secret = false

  supported_identity_providers = ["COGNITO"]
  explicit_auth_flows          = ["ALLOW_REFRESH_TOKEN_AUTH", "ALLOW_USER_SRP_AUTH"]

  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["openid", "profile"]

  callback_urls = concat(var.auth_callback_urls, ["https://${local.domain}"])
}

# ----------------------------------------------------------------------------------------------------------------------
# SSM Params
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_ssm_parameter" "cognito_pool_id" {
  name  = "/${var.application}/cognito/pool-id"
  type  = "String"
  value = aws_cognito_user_pool.main.id
}

resource "aws_ssm_parameter" "cognito_password_polic" {
  name  = "/${var.application}/cognito/password-policy"
  type  = "String"
  value = jsonencode(aws_cognito_user_pool.main.password_policy[0])
}

resource "aws_ssm_parameter" "cognito_client_dashboard" {
  name  = "/${var.application}/cognito/client-dashboard"
  type  = "String"
  value = aws_cognito_user_pool_client.dashboard.id
}