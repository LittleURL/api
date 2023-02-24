module "functions" {
  source = "./functions"

  prefix       = local.prefix
  lumigo_token = var.lumigo_token
  aws_account  = local.aws_account
  aws_region   = var.aws_region

  gateway_id            = aws_apigatewayv2_api.api.id
  gateway_execution_arn = aws_apigatewayv2_api.api.execution_arn
  authorizer_id         = aws_apigatewayv2_authorizer.cognito.id
  cognito_pool_id       = aws_cognito_user_pool.main.id
  cognito_pool_arn      = aws_cognito_user_pool.main.arn

  ddb_table_arns = {
    users      = aws_dynamodb_table.users.arn
    user_roles = aws_dynamodb_table.user_roles.arn
    domains    = aws_dynamodb_table.domains.arn
    links      = aws_dynamodb_table.links.arn
  }

  ddb_table_names = {
    users      = aws_dynamodb_table.users.id
    user_roles = aws_dynamodb_table.user_roles.id
    domains    = aws_dynamodb_table.domains.id
    links      = aws_dynamodb_table.links.id
  }
}
