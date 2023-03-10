locals {
  prefix            = "${local.application_clean}-"
  application_clean = lower(trimspace(var.application))
  environment       = contains(var.environments, terraform.workspace) ? terraform.workspace : "dev"
  aws_account       = lookup(var.aws_accounts, local.environment)
}

variable "application" {
  type        = string
  description = "Application name for prefixing globally unique resource names"
  default     = "LittleURL"

  validation {
    condition     = length(regex("^[\\w-]{1,64}$", var.application)) > 0
    error_message = "Application name value can only contain up to 64 alpha-numerical characters"
  }
}

variable "environments" {
  type    = set(string)
  default = ["dev", "prod"]
}

variable "lumigo_token" {
  type        = string
  description = "Lumigo API Token"
  default     = ""
}

variable "domain" {
  type = string
}

variable "email_spf_includes" {
  type    = list(string)
  default = []
}

# ----------------------------------------------------------------------------------------------------------------------
# AWS
# ----------------------------------------------------------------------------------------------------------------------
variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "aws_role" {
  type    = string
  default = "deploy-api"
}

variable "aws_accounts" {
  type = map(string)
  default = {
    dev  = "000000000000"
    prod = "000000000000"
  }
}

variable "aws_default_tags" {
  type        = map(string)
  description = "Common resource tags for all AWS resources"
  default = {
    application = "LittleURL"
  }
}

variable "aws_gateway_access_log" {
  type        = bool
  description = "Enable access logging for API Gateway"
  default     = false
}

# ----------------------------------------------------------------------------------------------------------------------
# API
# ----------------------------------------------------------------------------------------------------------------------
variable "cors_origins" {
  type = map(list(string))
  default = {
    dev  = ["*"]
    prod = ["https://littleurl.io"]
  }
}

variable "cors_expose" {
  type = map(list(string))
  default = {
    dev  = ["*"]
    prod = []
  }
}

variable "auth_callback_urls" {
  type        = list(string)
  description = "Additional callback domains for cognito client"
  default     = []
}
