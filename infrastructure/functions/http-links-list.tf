# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_links_list" {
  source = "../modules/lambda-function"

  aws_account = var.aws_account
  aws_region  = var.aws_region

  name          = "${var.prefix}http-links-list"
  source_key    = "http-links-list.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_links_list" {
  source = "../modules/lambda-gateway"
  method = "GET"
  path   = "/domains/{domainId}/links"

  function_name       = module.lambda_http_links_list.function_name
  function_invoke_arn = module.lambda_http_links_list.function_invoke_arn

  gateway_id            = var.gateway_id
  gateway_execution_arn = var.gateway_execution_arn

  authorizer_id = var.authorizer_id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_links_list_dynamodb" {
  source = "../modules/iam-dynamodb"
  role   = module.lambda_http_links_list.role_id

  tables = [
    {
      arn         = var.ddb_table_arns.user_roles
      enable_read = true
    },
    {
      arn         = var.ddb_table_arns.links
      enable_read = true
    }
  ]
}
