resource "aws_sqs_queue" "user_update" {
  name                      = "${local.prefix}user-update.fifo"
  fifo_queue                = true
  sqs_managed_sse_enabled   = true
  message_retention_seconds = 1209600 // 14days

  // dedup
  deduplication_scope   = "messageGroup"
  fifo_throughput_limit = "perMessageGroupId"
}
