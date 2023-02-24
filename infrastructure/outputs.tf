output "functions_bucket" {
  value = module.functions.functions_bucket
}

output "aws_assume_role" {
  value = "arn:aws:iam::${local.aws_account}:role/${var.aws_role}"
}