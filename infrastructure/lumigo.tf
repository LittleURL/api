locals {
  envvar_lumigo = var.lumigo_token == "" ? { "LUMIGO_USE_TRACER_EXTENSION" = false } : {
    "LUMIGO_USE_TRACER_EXTENSION" = true,
    "LUMIGO_TRACER_TOKEN"         = var.lumigo_token
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Role
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role" "lumigo_integration" {
  name               = "LumigoIntegration"
  assume_role_policy = data.aws_iam_policy_document.lumigo_integration_assume.json
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/ReadOnlyAccess",
    "arn:aws:iam::aws:policy/service-role/AWSAppSyncPushToCloudWatchLogs",
    "arn:aws:iam::aws:policy/AWSAppSyncAdministrator"
  ]
}

data "aws_iam_policy_document" "lumigo_integration_assume" {
  statement {
    actions = [
      "sts:AssumeRole"
    ]

    principals {
      type        = "AWS"
      identifiers = ["599735196807"]
    }
  }

  statement {
    actions = [
      "sts:AssumeRole"
    ]

    principals {
      type        = "AWS"
      identifiers = ["114300393969"]
    }
  }

  statement {
    actions = [
      "sts:AssumeRole"
    ]

    principals {
      type        = "Service"
      identifiers = ["appsync.amazonaws.com"]
    }
  }
}

# ----------------------------------------------------------------------------------------------------------------------
# Policies
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_iam_role_policy" "lumigo_integration" {
  name   = "LumigoIntegration"
  role   = aws_iam_role.lumigo_integration.id
  policy = data.aws_iam_policy_document.lumigo_integration.json
}

data "aws_iam_policy_document" "lumigo_integration" {
  statement {
    actions = [
      "lambda:UpdateFunctionConfiguration"
    ]
    resources = [
      "arn:aws:lambda:*:*:function:*"
    ]
  }

  statement {
    actions = [
      "ce:GetCostAndUsageWithResources",
      "ce:GetCostAndUsage"
    ]
    resources = [
      "*"
    ]
  }

  statement {
    actions = [
      "logs:PutSubscriptionFilter",
      "logs:DeleteSubscriptionFilter",
      "logs:DescribeSubscriptionFilters",
      "cloudwatch:PutMetricAlarm",
      "cloudwatch:DeleteAlarms",
      "cloudwatch:PutDashboard",
      "cloudwatch:DeleteDashboards",
      "cloudwatch:PutMetricData",
      "cloudwatch:PutMetricStream",
      "cloudwatch:DeleteMetricStream",
      "cloudwatch:StartMetricStreams",
      "cloudwatch:StopMetricStreams"
    ]
    resources = [
      "*"
    ]
  }
}

