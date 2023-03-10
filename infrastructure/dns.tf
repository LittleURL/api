resource "aws_route53_zone" "main" {
  name = var.domain
}

resource "aws_route53_record" "api" {
  zone_id = aws_route53_zone.main.zone_id
  type    = "A"
  name    = aws_apigatewayv2_domain_name.api.domain_name

  alias {
    name                   = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# SSM Params
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_ssm_parameter" "zone_id" {
  name  = "/${local.application_clean}/zone-id"
  type  = "String"
  value = aws_route53_zone.main.zone_id
}
