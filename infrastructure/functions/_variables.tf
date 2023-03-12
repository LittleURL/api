# ----------------------------------------------------------------------------------------------------------------------
# Misc
# ----------------------------------------------------------------------------------------------------------------------
variable "prefix" {
  type    = string
  default = ""
}

variable "environment" {
  type    = map(string)
  default = {}
}

variable "lumigo_token" {
  type    = string
  default = ""
}

variable "email_allowed_from_addresses" {
  type = list(string)
}

# ----------------------------------------------------------------------------------------------------------------------
# AWS
# ----------------------------------------------------------------------------------------------------------------------
variable "aws_region" {
  type = string
}

variable "aws_account" {
  type = string
}

variable "gateway_id" {
  type = string
}

variable "gateway_execution_arn" {
  type = string
}

variable "authorizer_id" {
  type = string
}

variable "cognito_pool_id" {
  type = string
}

variable "cognito_pool_arn" {
  type = string
}

variable "ddb_table_arns" {
  type = object({
    users        = string
    user_roles   = string
    user_invites = string
    domains      = string
    links        = string
  })
}

variable "ddb_table_names" {
  type = object({
    users        = string
    user_roles   = string
    user_invites = string
    domains      = string
    links        = string
  })
}
