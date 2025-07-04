### Create AWS foundational resources

# Create a Resource Group to identify Terraform created resources
resource "aws_resourcegroups_group" "Terraform" {
  name        = format("%s%s%s%s", var.Prefix, "rgg", var.EnvCode, "demoall")
  description = "Terraform created demo environment resources"

  resource_query {
    query = <<JSON
{
  "ResourceTypeFilters": [
    "AWS::AllSupported"
  ],
  "TagFilters": [
    {
      "Key": "Solution",
      "Values": ["${var.SolTag}"]
    }
  ]
}
JSON
  }

  tags = {
    Name  = format("%s%s%s%s", var.Prefix, "rgg", var.EnvCode, "demoall")
    rtype = "scaffold"
  }
}

# Create Amazon S3 bucket for ALB logs
resource "aws_s3_bucket" "alblogs" {
  bucket_prefix = format("%s%s%s%s", var.Prefix, "sss", var.EnvCode, "alblogs")
  force_destroy = true

  tags = {
    Name      = format("%s%s%s%s", var.Prefix, "sss", var.EnvCode, "alblogs"),
    rtype     = "storage"
    codeblock = "lzbase"
  }
}

# Create IAM Policy to enforce TLS 1.2 on Amazon S3 bucket and allow ALB access
data "aws_iam_policy_document" "S3logsTLS" {
  statement {
    sid    = "Allow HTTPS only"
    effect = "Deny"

    principals {
      type        = "*"
      identifiers = ["*"]
    }
    actions = [
      "s3*"
    ]
    resources = [
      "${aws_s3_bucket.alblogs.arn}",
      "${aws_s3_bucket.alblogs.arn}/*"
    ]
    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values = [
        "false"
      ]
    }
  }
  statement {
    sid    = "Allow TLS 1.2 and above"
    effect = "Deny"

    principals {
      type        = "*"
      identifiers = ["*"]
    }
    actions = [
      "s3*"
    ]
    resources = [
      "${aws_s3_bucket.alblogs.arn}",
      "${aws_s3_bucket.alblogs.arn}/*"
    ]
    condition {
      test     = "NumericLessThan"
      variable = "s3:TlsVersion"
      values = [
        "1.2"
      ]
    }
  }
  statement {
    # https://docs.aws.amazon.com/elasticloadbalancing/latest/application/enable-access-logging.html
    # Consider limiting principle to specific region
    sid    = "Allow ALB logging access regions available as of August 2022 or later"
    effect = "Allow"

    principals {
      type = "AWS"
      identifiers = ["arn:aws:iam::127311923021:root", #US East (N. Virginia) 
        "arn:aws:iam::033677994240:root",              #US East (Ohio)
        "arn:aws:iam::027434742980:root",              #US West (N. California)
        "arn:aws:iam::797873946194:root",              #US West (Oregon)
        "arn:aws:iam::098369216593:root",              #Africa (Cape Town)
        "arn:aws:iam::754344448648:root",              #Asia Pacific (Hong Kong)
        "arn:aws:iam::589379963580:root",              #Asia Pacific (Jakarta)
        "arn:aws:iam::718504428378:root",              #Asia Pacific (Mumbai)
        "arn:aws:iam::383597477331:root",              #Asia Pacific (Osaka)
        "arn:aws:iam::600734575887:root",              #Asia Pacific (Seoul)
        "arn:aws:iam::114774131450:root",              #Asia Pacific (Singapore)
        "arn:aws:iam::783225319266:root",              #Asia Pacific (Sydney)
        "arn:aws:iam::582318560864:root",              #Asia Pacific (Tokyo) 
        "arn:aws:iam::985666609251:root",              #Canada (Central) 
        "arn:aws:iam::054676820928:root",              #Europe (Frankfurt)
        "arn:aws:iam::156460612806:root",              #Europe (Ireland)
        "arn:aws:iam::652711504416:root",              #Europe (London)
        "arn:aws:iam::635631232127:root",              #Europe (Milan)
        "arn:aws:iam::009996457667:root",              #Europe (Paris) 
        "arn:aws:iam::897822967062:root",              #Europe (Stockholm)
        "arn:aws:iam::076674570225:root",              #Middle East (Bahrain)
      "arn:aws:iam::507241528517:root"]                #South America (São Paulo) 
    }
    actions = [
      "s3:PutObject"
    ]
    resources = [
      "${aws_s3_bucket.alblogs.arn}",
      "${aws_s3_bucket.alblogs.arn}/*"
    ]
  }
}

# Apply policy to enforce TLS 1.2 on Amazon S3 buckets
resource "aws_s3_bucket_policy" "alblogs" {
  bucket = aws_s3_bucket.alblogs.id
  policy = data.aws_iam_policy_document.S3logsTLS.json
}

# Enable Amazon S3 Bucket versioning
resource "aws_s3_bucket_versioning" "alblogs" {
  bucket = aws_s3_bucket.alblogs.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Block Amazon S3 Bucket public access
resource "aws_s3_bucket_public_access_block" "alblogs" {
  bucket                  = aws_s3_bucket.alblogs.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}