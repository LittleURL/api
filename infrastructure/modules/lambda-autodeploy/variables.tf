variable "name" {
  type = string
}

variable "role_path" {
  type    = string
  default = "/lambda/"
}

variable "bucket" {
  type = string
}

variable "aws_region" {
  type = string
}

variable "aws_account" {
  type = string
}

variable "function_name_prefix" {
  type = string
  default = ""
}