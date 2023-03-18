# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_links_create" {
  source = "../modules/lambda-function"

  aws_account       = var.aws_account
  aws_region        = var.aws_region
  enable_autodeploy = var.enable_autodeploy

  name          = "${var.prefix}http-links-create"
  source_key    = "http-links-create.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_links_create" {
  source = "../modules/lambda-gateway"
  method = "POST"
  path   = "/domains/{domainId}/links"

  function_name       = module.lambda_http_links_create.function_name
  function_invoke_arn = module.lambda_http_links_create.function_invoke_arn

  gateway_id            = var.gateway_id
  gateway_execution_arn = var.gateway_execution_arn

  authorizer_id = var.authorizer_id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_links_create_dynamodb" {
  source = "../modules/iam-dynamodb"
  role   = module.lambda_http_links_create.role_id

  tables = [
    {
      arn         = var.ddb_table_arns.user_roles
      enable_read = true
    },
    {
      arn         = var.ddb_table_arns.links
      enable_write = true
    }
  ]
}
