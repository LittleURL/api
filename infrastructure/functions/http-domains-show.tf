# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_domains_show" {
  source = "../modules/lambda-function"

  aws_account = var.aws_account
  aws_region  = var.aws_region

  name          = "${var.prefix}http-domains-show"
  source_key    = "http-domains-show.zip"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_domains_show" {
  source = "../modules/lambda-gateway"
  method = "GET"
  path   = "/domains/{domainId}"

  function_name       = module.lambda_http_domains_show.function_name
  function_invoke_arn = module.lambda_http_domains_show.function_invoke_arn

  gateway_id            = var.gateway_id
  gateway_execution_arn = var.gateway_execution_arn

  authorizer_id = var.authorizer_id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_domains_show_dynamodb" {
  source = "../modules/iam-dynamodb"
  role   = module.lambda_http_domains_show.role_id

  tables = [
    {
      arn         = var.ddb_table_arns.domains
      enable_read = true
    },
    {
      arn         = var.ddb_table_arns.user_roles
      enable_read = true
    }
  ]
}
