# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_domains_update" {
  source = "./modules/lambda-function"

  aws_account = local.aws_account
  aws_region  = var.aws_region

  name          = "${local.prefix}http-domains-update"
  source_key    = "http-domains-update.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_domains_update" {
  source = "./modules/lambda-gateway"
  method = "PATCH"
  path   = "/companies/{domainId}"

  function_name       = module.lambda_http_domains_update.function_name
  function_invoke_arn = module.lambda_http_domains_update.function_invoke_arn

  gateway_id            = aws_apigatewayv2_api.api.id
  gateway_execution_arn = aws_apigatewayv2_api.api.execution_arn

  authorizer_id = aws_apigatewayv2_authorizer.auth0.id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_domains_update_dynamodb" {
  source = "./modules/iam-dynamodb"
  role   = module.lambda_http_domains_update.role_id
  table  = aws_dynamodb_table.domains.arn

  enable_write = true
}
