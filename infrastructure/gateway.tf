locals {
  domain_api = "api.${aws_route53_zone.main.name}"
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

  dynamic "access_log_settings" {
    for_each = var.aws_gateway_access_log ? ["enabled"] : []

    content {
      destination_arn = aws_cloudwatch_log_group.gateway_v1[0].arn
      format = jsonencode({
        httpMethod     = "$context.httpMethod",
        ip             = "$context.identity.sourceIp",
        protocol       = "$context.protocol",
        requestId      = "$context.requestId",
        requestTime    = "$context.requestTime",
        responseLength = "$context.responseLength",
        routeKey       = "$context.routeKey",
        status         = "$context.status",
        error = {
          message      = "$context.error.message"
          responseType = "$context.error.responseType",
        },
        authorizer = {
          error       = "$context.authorizer.error",
          principalId = "$context.authorizer.principalId",
        },
        integration = {
          status            = "$context.integration.status",
          error             = "$context.integration.error",
          integrationStatus = "$context.integration.integrationStatus",
          latency           = "$context.integration.latency",
          requestId         = "$context.integration.requestId",
        }
      })
    }
  }

  lifecycle {
    // auto-deploy changes this
    ignore_changes = [deployment_id]
  }
}

resource "aws_cloudwatch_log_group" "gateway_v1" {
  count = var.aws_gateway_access_log ? 1 : 0

  name = "APIGatewayV2/${local.prefix}api/v1"
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

  depends_on = [
    aws_acm_certificate.api,
    aws_route53_record.api_cert_validation[0],
    aws_route53_record.api_cert_validation[1],
    aws_route53_record.api_cert_validation[2],
  ]
}

resource "aws_apigatewayv2_api_mapping" "api" {
  api_id          = aws_apigatewayv2_api.api.id
  domain_name     = aws_apigatewayv2_domain_name.api.id
  stage           = aws_apigatewayv2_stage.v1.id
  api_mapping_key = "v1"
}

resource "aws_acm_certificate" "api" {
  domain_name       = local.domain_api
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "api_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.api.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id         = aws_route53_zone.main.zone_id
  name            = each.value.name
  type            = each.value.type
  ttl             = 60
  records         = [each.value.record]
  allow_overwrite = true
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
    audience = [aws_cognito_user_pool_client.dashboard.id]
  }
}
