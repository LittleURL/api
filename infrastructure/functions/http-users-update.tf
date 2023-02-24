# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_users_update" {
  source = "../modules/lambda-function"

  aws_account = var.aws_account
  aws_region  = var.aws_region

  name          = "${var.prefix}http-users-update"
  source_key    = "http-users-update.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_users_update" {
  source = "../modules/lambda-gateway"
  method = "PUT"
  path   = "/domains/{domainId}/users/{userId}"

  function_name       = module.lambda_http_users_update.function_name
  function_invoke_arn = module.lambda_http_users_update.function_invoke_arn

  gateway_id            = var.gateway_id
  gateway_execution_arn = var.gateway_execution_arn

  authorizer_id = var.authorizer_id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_users_update_dynamodb" {
  source = "../modules/iam-dynamodb"
  role   = module.lambda_http_users_update.role_id

  tables = [
    {
      arn           = var.ddb_table_arns.user_roles
      enable_read   = true
      enable_write  = true
    }
  ]
}
