resource "aws_iam_role_policy" "dynamodb" {
  name   = var.name
  role   = var.role
  policy = data.aws_iam_policy_document.dynamodb.json
}

data "aws_iam_policy_document" "dynamodb" {
  // Read
  dynamic "statement" {
    for_each = var.enable_read == true ? ["Read"] : []
    content {
      sid = "Read"
      actions = [
        "dynamodb:GetItem",
        "dynamodb:Scan",
        "dynamodb:Query",
        "dynamodb:BatchGetItem"
      ]
      resources = [
        var.table,
        "${var.table}/index/*"
      ]
    }
  }

  // Write
  dynamic "statement" {
    for_each = var.enable_write == true ? ["Write"] : []
    content {
      sid = "Write"
      actions = [
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:BatchWriteItem",
        "dynamodb:UpdateTimeToLive"
      ]
      resources = [var.table]
    }
  }

  // Delete
  dynamic "statement" {
    for_each = var.enable_delete == true ? ["Delete"] : []
    content {
      sid = "Delete"
      actions = [
        "dynamodb:Deleteitem"
      ]
      resources = [var.table]
    }
  }

  // Stream
  dynamic "statement" {
    for_each = var.enable_stream == true ? ["Stream"] : []
    content {
      sid = "Stream"
      actions = [
        "dynamodb:ListStreams",
        "dynamodb:DescribeStream",
        "dynamodb:GetRecords",
        "dynamodb:GetShardIterator"
      ]
      resources = [
        "${var.table}/stream/*"
      ]
    }
  }
}
