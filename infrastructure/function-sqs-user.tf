# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
module "lambda_sqs_user_update" {
  source = "./modules/lambda-function"

  aws_account = local.aws_account
  aws_region  = var.aws_region

  prefix        = local.prefix
  name          = "sqs-user-update"
  source_bucket = aws_s3_bucket.functions.id

  environment_variables = merge(local.envvar_default, {
    AUTH0_DOMAIN       = var.auth0_domain,
    AUTH0_CLIENTID     = auth0_client.api.client_id,
    AUTH0_CLIENTSECRET = auth0_client.api.client_secret
  })
}

module "lambda_sqs_user_update_event_source" {
  source = "./modules/lambda-sqs"

  function_name     = module.lambda_sqs_user_update.function_name
  function_role_arn = module.lambda_sqs_user_update.role_arn
  queue_arn         = aws_sqs_queue.user_update.arn
  batch_size        = 1
}
