# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_http_domains_list" {
  source = "./modules/lambda-function"

  aws_account = local.aws_account
  aws_region  = var.aws_region

  prefix        = local.prefix
  name          = "http-domains-list"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {})
}

module "gateway_lambda_http_domains_list" {
  source = "./modules/lambda-gateway"
  method = "POST"
  path   = "/companies"

  function_name       = module.lambda_http_domains_list.function_name
  function_invoke_arn = module.lambda_http_domains_list.function_invoke_arn

  gateway_id            = aws_apigatewayv2_api.api.id
  gateway_execution_arn = aws_apigatewayv2_api.api.execution_arn

  authorizer_id = aws_apigatewayv2_authorizer.auth0.id
}

# ----------------------------------------------------------------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role_policy" "lambda_http_domains_list_dynamo" {
  name   = "DynamoDB"
  role   = module.lambda_http_domains_list.role_id
  policy = data.aws_iam_policy_document.lambda_http_domains_list_dynamo.json
}

data "aws_iam_policy_document" "lambda_http_domains_list_dynamo" {
  statement {
    sid = "EventSourceMapping"
    actions = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes"
    ]
    resources = [
      aws_sqs_queue.user_update.arn
    ]
  }
}