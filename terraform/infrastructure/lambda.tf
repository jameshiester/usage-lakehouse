module "go_lambda_function" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = format("%s-%s-%s", var.Prefix, var.EnvCode, "migration")

  attach_cloudwatch_logs_policy     = false
  cloudwatch_logs_retention_in_days = 1
  timeout                           = 120
  tags                              = local.tags
  handler                           = "bootstrap"
  runtime                           = "provided.al2023"
  architectures                     = ["arm64"] # x86_64 (GOARCH=amd64); arm64 (GOARCH=arm64)

  trigger_on_package_timestamp = true

  vpc_subnet_ids         = module.vpc.private_subnets
  vpc_security_group_ids = [aws_security_group.lambda.id]
  attach_network_policy  = true
  attach_policy_json     = true
  logging_log_group      = format("%s-%s-%s", var.Prefix, var.EnvCode, "migration")
  environment_variables = {
    LAMBDA_MODE          = "true"
    POSTGRES_DB          = module.db.db_instance_database
    POSTGRES_HOST        = module.db.db_instance_endpoint
    POSTGRES_USERNAME    = module.db.db_instance_username
    POSTGRES_PORT        = module.db.db_instance_port
    DB_MASTER_SECRET_ARN = module.db.db_instance_master_user_secret_arn
  }
  policy_json = jsonencode(
    {
      Version = "2012-10-17",
      Statement = [
        {
          Effect = "Allow",
          Action = ["secretsmanager:DescribeSecret",
            "secretsmanager:GetSecretValue"
          ]
          Resource = [module.db.db_instance_master_user_secret_arn]
        },
        {
          Effect : "Allow",
          Action : [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents",
          ],
          Resource : "*"
        }
      ]
    }
  )

  source_path = [
    {
      path = "${path.module}/../../go"
      commands = [
        "cp -r cmd/migrations ../../",
        "rm main.go",
        "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap cmd/migrations/main.go",
        ":zip",
      ]
      patterns = [
        "!cmd/migrations/main.go",
        "cmd/migrations/.*",
        "bootstrap"
     
      ]
    }
  ]
}

resource "aws_lambda_invocation" "example" {
  function_name = module.go_lambda_function.lambda_function_name
  input = jsonencode({
    args = ["up"]
  })
  triggers = {
    hash = module.go_lambda_function.lambda_function_source_code_hash
  }
}