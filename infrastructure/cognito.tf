locals {
  cognito_domain = "auth.${local.domain}"
}

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
  domain          = local.cognito_domain
  user_pool_id    = aws_cognito_user_pool.main.id
  certificate_arn = aws_acm_certificate.cognito.arn
}

resource "cloudflare_record" "cognito" {
  zone_id = local.zone_id
  name    = local.cognito_domain
  type    = "CNAME"
  value   = aws_cognito_user_pool_domain.main.cloudfront_distribution_arn
  ttl     = 600
}

resource "aws_acm_certificate" "cognito" {
  domain_name       = local.cognito_domain
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "cloudflare_record" "cognito_cert" {
  for_each = {
    for dvo in aws_acm_certificate.cognito.domain_validation_options : dvo.domain_name => {
      resource_record_name  = dvo.resource_record_name
      resource_record_value = dvo.resource_record_value
      resource_record_type  = dvo.resource_record_type
    }
  }

  zone_id = local.zone_id
  proxied = false

  name  = each.value.resource_record_name
  type  = each.value.resource_record_type
  value = trimsuffix(each.value.resource_record_value, ".")
}

resource "aws_acm_certificate_validation" "cognito" {
  certificate_arn         = aws_acm_certificate.cognito.arn
  validation_record_fqdns = [for record in cloudflare_record.cognito_cert : record.hostname]
}
