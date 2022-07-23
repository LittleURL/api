variable "name" {
  type        = string
  description = "Name of the IAM policy"
  default     = "DynamoDB"
}

variable "role" {
  type        = string
  description = "ARN of IAM role to attach the policy to"
}

variable "table" {
  type        = string
  description = "ARN of DynamoDB table"
}

variable "enable_read" {
  type        = bool
  description = "Allow reading of items from the table"
  default     = false
}

variable "enable_write" {
  type        = bool
  description = "Allow writing items to the table"
  default     = false
}

variable "enable_delete" {
  type        = bool
  description = "Allow deletion of items from the table"
  default     = false
}

variable "enable_stream" {
  type        = bool
  description = "Allow access to DynamoDB Stream"
  default     = false
}
