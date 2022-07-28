resource "aws_cognito_user_pool" "main" {
  name = var.application

  mfa_configuration = "OPTIONAL"

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
    pre_token_generation = module.lambda_cognito_pre_token_generation.function_arn
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Domain
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_cognito_user_pool_domain" "main" {
  domain       = local.environment == "prod" ? var.application : "${local.prefix}${terraform.workspace}"
  user_pool_id = aws_cognito_user_pool.main.id
}

resource "cloudflare_record" "cognito" {
  zone_id = local.zone_id
  name    = "auth.${local.domain}"
  type    = "CNAME"
  value   = aws_cognito_user_pool_domain.main.cloudfront_distribution_arn
  ttl     = 600
}
