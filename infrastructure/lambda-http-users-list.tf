# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_users_list" {
  source = "./modules/lambda-function"

  aws_account = local.aws_account
  aws_region  = var.aws_region

  name          = "${local.prefix}http-users-list"
  source_key    = "http-users-list.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_users_list" {
  source = "./modules/lambda-gateway"
  method = "GET"
  path   = "/domains/{domainId}/users"

  function_name       = module.lambda_http_users_list.function_name
  function_invoke_arn = module.lambda_http_users_list.function_invoke_arn

  gateway_id            = aws_apigatewayv2_api.api.id
  gateway_execution_arn = aws_apigatewayv2_api.api.execution_arn

  authorizer_id = aws_apigatewayv2_authorizer.cognito.id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_users_list_dynamodb" {
  source = "./modules/iam-dynamodb"
  role   = module.lambda_http_users_list.role_id

  tables = [
    {
      arn         = aws_dynamodb_table.user_roles.arn
      enable_read = true
    },
    {
      arn         = aws_dynamodb_table.users.arn
      enable_read = true
    }
  ]
}
