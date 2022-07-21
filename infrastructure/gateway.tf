# ----------------------------------------------------------------------------------------------------------------------
# Gateway
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_apigatewayv2_api" "api" {
  name          = "public-companies"
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
    certificate_arn = aws_acm_certificate.api.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "cloudflare_record" "api" {
  zone_id = cloudflare_zone.main.id
  name    = "api"
  type    = "CNAME"
  value   = aws_apigatewayv2_domain_name.api.domain_name_configuration.0.target_domain_name
  proxied = true
}

# ----------------------------------------------------------------------------------------------------------------------
# Origin Cert
# ----------------------------------------------------------------------------------------------------------------------
resource "tls_private_key" "api" {
  algorithm = "RSA"
}

resource "tls_cert_request" "api" {
  private_key_pem = tls_private_key.api.private_key_pem

  dns_names = [local.domain_api]

  subject {
    common_name  = local.domain_api
    organization = "LittleURL"
  }
}

resource "cloudflare_origin_ca_certificate" "api" {
  csr                = tls_cert_request.api.cert_request_pem
  hostnames          = [local.domain_api]
  request_type       = "origin-rsa"
  requested_validity = 5475 // (15yrs) Cloudflare default
}

data "http" "cloudflare_origin_root_ca" {
  url = "https://developers.cloudflare.com/ssl/static/origin_ca_rsa_root.pem"
}

resource "aws_acm_certificate" "api" {
  private_key       = tls_private_key.api.private_key_pem
  certificate_body  = cloudflare_origin_ca_certificate.api.certificate
  certificate_chain = data.http.cloudflare_origin_root_ca.response_body
}

# ----------------------------------------------------------------------------------------------------------------------
# Authorizer
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_apigatewayv2_authorizer" "auth0" {
  api_id           = aws_apigatewayv2_api.api.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "auth0"

  jwt_configuration {
    issuer   = "https://${var.auth0_domain}"
    audience = ["https://${local.domain_api}"]
  }
}
