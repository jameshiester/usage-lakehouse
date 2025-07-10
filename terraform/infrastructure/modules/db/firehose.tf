data "aws_iam_policy_document" "firehose_assume_role" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["firehose.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "firehose" {
  statement {
    # https://docs.aws.amazon.com/kms/latest/developerguide/key-policy-overview.html
    sid    = "Enable IAM User Permissions"
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
    actions = [
      "kms*"
    ]
    resources = [
      "*"
    ]
  }
  statement {
    # https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/encrypt-log-data-kms.html
    sid    = "Allow Cloudwatch access to KMS Key"
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["logs.${var.Region}.amazonaws.com"]
    }
    actions = [
      "kms:Encrypt*",
      "kms:Decrypt*",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:Describe*"
    ]
    resources = [
      "*"
    ]
    condition {
      test     = "ArnLike"
      variable = "kms:EncryptionContext:aws:logs:arn"
      values = [
        "arn:aws:logs:${var.Region}:${data.aws_caller_identity.current.account_id}:*"
      ]
    }
  }
}

# Create KMS key for solution
resource "aws_kms_key" "firehose" {
  description             = "KMS key to secure firehose"
  deletion_window_in_days = 7
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.firehose.json

  tags = {
    Name         = format("%s%s%s%s", var.Prefix, "kms", var.EnvCode, "firehose")
    resourcetype = "security"
    codeblock    = "ecscluster"
  }
}

# Create KMS Alias. Only used in this context to provide a friendly display name
resource "aws_kms_alias" "firehose" {
  name          = "alias/firehose"
  target_key_id = aws_kms_key.firehose.key_id
}

data "aws_iam_policy_document" "firehose_s3" {
  statement {
    effect = "Allow"
    actions = [
      "glue:GetTable",
      "glue:GetDatabase",
      "glue:UpdateTable",
      "glue:CreateTable",
      "glue:CreateDatabase"
    ]
    resources = [
      "arn:aws:glue:${var.Region}:${data.aws_caller_identity.current.account_id}:catalog",
      "arn:aws:glue:${var.Region}:${data.aws_caller_identity.current.account_id}:database/*",
      "arn:aws:glue:${var.Region}:${data.aws_caller_identity.current.account_id}:table/*/*"
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:AbortMultipartUpload",
      "s3:GetBucketLocation",
      "s3:GetObject",
      "s3:ListBucket",
      "s3:ListBucketMultipartUploads",
      "s3:PutObject",
      "s3:DeleteObject"
    ]
    resources = [
      aws_s3_bucket.storage.arn,
      "${aws_s3_bucket.storage.arn}/*"
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "kms:Decrypt",
      "kms:GenerateDataKey"
    ]
    resources = [
      aws_kms_key.firehose.arn
    ]
    condition {
      test     = "StringEquals"
      variable = "kms:ViaService"

      values = [
        "s3.region.amazonaws.com",
      ]
    }
    condition {
      test     = "StringLike"
      variable = "kms:EncryptionContext:aws:s3:arn"

      values = [
        "${aws_s3_bucket.storage.arn}/prefix*",
      ]
    }

  }
  statement {
    effect = "Allow"
    actions = [
      "logs:PutLogEvents"
    ]
    resources = [
      "*"
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "secretsmanager:GetSecretValue"
    ]
    resources = [
      module.db.db_instance_master_user_secret_arn
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "ec2:DescribeVpcEndpointServices"
    ]
    resources = [
      "*"
    ]
  }
}

# Create Amazon S3 bucket for ALB logs
resource "aws_s3_bucket" "storage" {
  bucket_prefix = format("%s-%s-%s", var.Prefix, "firehose", var.EnvCode)
  force_destroy = true

  tags = local.tags
}

# Create IAM Role for Firehose
resource "awscc_iam_role" "firehose" {
  role_name                   = format("%s-%s-%s", var.Prefix, "firehose", var.EnvCode)
  assume_role_policy_document = jsonencode(jsondecode(data.aws_iam_policy_document.firehose_assume_role.json))
  managed_policy_arns         = []
  path                        = "/"
  tags = local.tags
}

resource "awscc_iam_role_policy" "firehose" {
  policy_document = jsonencode(jsondecode(data.aws_iam_policy_document.firehose_s3.json))
  policy_name     = "firehose-s3-access"
  role_name       = awscc_iam_role.firehose.role_name
}

# Create the Kinesis Firehose Delivery Stream
resource "awscc_kinesisfirehose_delivery_stream" "example" {
  delivery_stream_name = format("%s-%s-%s", var.Prefix, "usage-lakehouse", var.EnvCode)
  iceberg_destination_configuration = {
    append_only = false
    s3_configuration = {
      bucket_arn = aws_s3_bucket.storage.arn
    }
    role_arn            = awscc_iam_role.firehose.arn
  }
  database_source_configuration = {
    columns = {
        include = ["*"]
    }
    tables = {
        include = ["account"]
    }
    database_source_authentication_configuration = {
      secrets_manager_configuration = {
        enabled = true
        role_arn = awscc_iam_role.firehose.arn
        secret_arn = module.db.db_instance_domain_auth_secret_arn
        
      }
      vpc_endpoint_service_name = aws_vpc_endpoint_service.rds_lb_endpoint_service.name
    }
    databases = {
        include = [var.DBInstanceDatabaseName]
    }

  }
  tags = [{key = "Envrionment",value= var.EnvTag}]
}