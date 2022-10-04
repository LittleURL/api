output "functions_bucket" {
  value = aws_s3_bucket.functions.id
}

output "aws_assume_role" {
  value = "arn:aws:iam::${local.aws_account}:role/${var.aws_role}"
}