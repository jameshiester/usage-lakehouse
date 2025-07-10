resource "aws_lb_target_group" "rds_target_group" {
  name        = format("%s-%s-%s", var.Prefix, "db-listener", var.EnvCode)
  port        = module.db.db_instance_port
  protocol    = "TCP"
  vpc_id      = var.VPCID
  target_type = "ip"
}

data "dns_a_record_set" "rds_ip" {
  host = module.db.db_instance_address
}

data "aws_iam_policy_document" "lambda_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/lambda_function.py"
  output_path = "${path.module}/lambda_function.zip"
}

# Attach a target to each target group
resource "aws_lb_target_group_attachment" "rds_target_group_attachment" {

  target_group_arn = aws_lb_target_group.rds_target_group.arn
  target_id        = data.dns_a_record_set.rds_ip.addrs[0]

  lifecycle {
    ignore_changes = [target_id]
  }
  depends_on = [aws_lb_target_group.rds_target_group]
}

data "aws_subnet" "selected" {
  filter {
    name   = "tag:Name"
    values = [var.VPCDatabaseSubnetGroup]
  }
}

resource "aws_lb" "rds_lb" {
  name                             = format("%s-%s-%s", var.Prefix, "db-lb", var.EnvCode)
  internal                         = true
  load_balancer_type               = "network"
  subnets                          = [data.aws_subnet.selected.id]
  enable_cross_zone_load_balancing = true
  tags                             = local.tags
}

# Create listeners for each RDS instance, mapping each to its respective target group
resource "aws_lb_listener" "rds_listener" {

  load_balancer_arn = aws_lb.rds_lb.arn
  port              = module.db.db_instance_port
  protocol          = "TCP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.rds_target_group.arn
  }
}

resource "aws_sns_topic" "default" {
  name = format("%s-%s-%s", var.Prefix, "rds-event", var.EnvCode)
}

resource "aws_db_event_subscription" "default" {
  name      = format("%s-%s-%s", var.Prefix, "rds-event-sub", var.EnvCode)
  sns_topic = aws_sns_topic.default.arn

  source_type = "db-instance"
  source_ids  = [module.db.db_instance_identifier]

  event_categories = [
    "failover",
    "failure"
  ]
}

# Create an IAM policy for the Lambda function
resource "aws_iam_role" "lambda_execution_role" {
  name               = format("%s-%s-%s", var.Prefix, "rds-event-sub-lambda", var.EnvCode)
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role_policy.json
}

# Create a Lambda function to check the RDS instance IP address
resource "aws_lambda_function" "check_rds_ip" {
  function_name = format("%s-%s-%s", var.Prefix, "update-rds-ip", var.EnvCode)
  role          = aws_iam_role.lambda_execution_role.arn
  handler       = "lambda_function.lambda_handler"
  runtime       = "python3.11"

  filename = data.archive_file.lambda_zip.output_path

  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  environment {
    variables = {
      Cluster_EndPoint = module.db.db_instance_endpoint
      RDS_Port : module.db.db_instance_port
      NLB_TG_ARN : aws_lb_target_group.rds_target_group.arn
    }
  }
}

# Create an IAM policy for the Lambda function
resource "aws_iam_role_policy" "lambda_execution_role_policy" {
  name   = format("%s-%s-%s", var.Prefix, "rds-ip-update-lambda", var.EnvCode)
  role   = aws_iam_role.lambda_execution_role.id
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "rds:DescribeDBInstances",
        "elasticloadbalancing:DescribeTargetHealth",
        "elasticloadbalancing:RegisterTargets",
        "elasticloadbalancing:DeregisterTargets"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_sns_topic_policy" "default" {
  arn = aws_sns_topic.default.arn

  policy = data.aws_iam_policy_document.sns_topic_policy.json
}

data "aws_iam_policy_document" "sns_topic_policy" {
  policy_id = "__default_policy_ID"

  statement {
    actions = [
      "SNS:GetTopicAttributes",
      "SNS:SetTopicAttributes",
      "SNS:AddPermission",
      "SNS:RemovePermission",
      "SNS:DeleteTopic",
      "SNS:Subscribe",
      "SNS:ListSubscriptionsByTopic",
      "SNS:Publish",
      "SNS:Receive"
    ]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceOwner"

      values = [
        data.aws_caller_identity.current.account_id,
      ]
    }

    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    resources = [
      aws_sns_topic.default.arn,
    ]

    sid = "__default_statement_ID"
  }
}


resource "aws_sns_topic_subscription" "user_updates_sqs_target" {
  topic_arn = aws_sns_topic.default.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.check_rds_ip.arn
}

# Create VPC endpoint service for the Load Balancer
resource "aws_vpc_endpoint_service" "rds_lb_endpoint_service" {
  acceptance_required        = false
  network_load_balancer_arns = [aws_lb.rds_lb.arn]

  supported_regions = [var.Region]

  tags = {
    Name = format("%s-%s-%s", var.Prefix, "rds-endpoint", var.EnvCode)
  }
}