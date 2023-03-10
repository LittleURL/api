# ----------------------------------------------------------------------------------------------------------------------
# Common Variables
# ----------------------------------------------------------------------------------------------------------------------
variable "aws_region" {
  type = string
}

variable "aws_account" {
  type = string
}

# ----------------------------------------------------------------------------------------------------------------------
# Function
# ----------------------------------------------------------------------------------------------------------------------
variable "name" {
  type = string
}

variable "role_path" {
  type    = string
  default = "/lambda/"
}

variable "environment_variables" {
  type    = map(string)
  default = {}
}

variable "source_bucket" {
  type = string
}

variable "source_key" {
  type        = string
  description = "Path of the file within specific S3 bucket, defaults to `{name}.zip`"
  default     = null
  nullable    = true
}

# ----------------------------------------------------------------------------------------------------------------------
# Runtime
# ----------------------------------------------------------------------------------------------------------------------
variable "architectures" {
  type    = set(string)
  default = ["arm64"]
}

variable "timeout" {
  type    = number
  default = 5
}

variable "memory" {
  type    = number
  default = 128
}
