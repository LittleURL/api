locals {
  domain_api = "api.${local.domain}"
}

data "aws_ssm_parameter" "cloudflare_zone" {
  name = "/${var.application}/cloudflare-zone"
}

data "cloudflare_zone" "default" {
  zone_id = data.aws_ssm_parameter.cloudflare_zone.value
}

data "aws_ssm_parameter" "api_origin_cert" {
  name = "/${var.application}/api-certificate-arn"
}

# ----------------------------------------------------------------------------------------------------------------------
# Gateway
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_apigatewayv2_api" "api" {
  name          = "${local.prefix}api"
  protocol_type = "HTTP"

  disable_execute_api_endpoint = true

  cors_configuration {
    allow_headers  = ["*"]
    allow_methods  = ["*"]
    allow_origins  = lookup(var.cors_origins, local.environment)
    expose_headers = lookup(var.cors_expose, local.environment)
  }
}

resource "aws_apigatewayv2_stage" "v1" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "v1"
  auto_deploy = true

  // required due to bug https://github.com/hashicorp/terraform-provider-aws/issues/14742#issuecomment-750693332
  default_route_settings {
    throttling_burst_limit = 100
    throttling_rate_limit  = 50
  }

  lifecycle {
    // auto-deploy changes this
    ignore_changes = [deployment_id]
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Domain
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_apigatewayv2_domain_name" "api" {
  domain_name = local.domain_api

  domain_name_configuration {
    certificate_arn = data.aws_ssm_parameter.api_origin_cert.value
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "api" {
  api_id          = aws_apigatewayv2_api.api.id
  domain_name     = aws_apigatewayv2_domain_name.api.id
  stage           = aws_apigatewayv2_stage.v1.id
  api_mapping_key = "v1"
}

resource "cloudflare_record" "api" {
  zone_id = local.zone_id
  name    = "api"
  type    = "CNAME"
  value   = aws_apigatewayv2_domain_name.api.domain_name_configuration.0.target_domain_name
  proxied = true
}

# ----------------------------------------------------------------------------------------------------------------------
# Authorizer
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id          = aws_apigatewayv2_api.api.id
  name            = "cognito"
  authorizer_type = "JWT"

  identity_sources = [
    "$request.header.Authorization"
  ]

  jwt_configuration {
    issuer   = "https://${aws_cognito_user_pool.main.endpoint}"
    audience = [local.domain_api]
  }
}
