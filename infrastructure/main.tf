module "functions" {
  source = "./functions"

  prefix            = local.prefix
  lumigo_token      = var.lumigo_token
  aws_account       = local.aws_account
  aws_region        = var.aws_region
  enable_autodeploy = var.function_autodeploy

  gateway_id                   = aws_apigatewayv2_api.api.id
  gateway_execution_arn        = aws_apigatewayv2_api.api.execution_arn
  authorizer_id                = aws_apigatewayv2_authorizer.cognito.id
  cognito_pool_id              = aws_cognito_user_pool.main.id
  cognito_pool_arn             = aws_cognito_user_pool.main.arn
  email_allowed_from_addresses = [local.email_from, local.email_from_friendly]

  environment = {
    AppName         = var.application
    EmailFrom       = local.email_from_friendly
    DashboardDomain = local.dashboard_domain
  }

  ddb_table_arns = {
    users        = aws_dynamodb_table.users.arn
    user_roles   = aws_dynamodb_table.user_roles.arn
    user_invites = aws_dynamodb_table.user_invites.arn
    domains      = aws_dynamodb_table.domains.arn
    links        = aws_dynamodb_table.links.arn
  }

  ddb_table_names = {
    users        = aws_dynamodb_table.users.id
    user_roles   = aws_dynamodb_table.user_roles.id
    user_invites = aws_dynamodb_table.user_invites.id
    domains      = aws_dynamodb_table.domains.id
    links        = aws_dynamodb_table.links.id
  }
}

module "functions_autodeploy" {
  count  = var.function_autodeploy ? 1 : 0
  source = "./modules/lambda-autodeploy"

  name                 = "${local.prefix}s3-autodeploy"
  bucket               = module.functions.functions_bucket
  aws_account          = local.aws_account
  aws_region           = var.aws_region
  function_name_prefix = local.prefix
}
