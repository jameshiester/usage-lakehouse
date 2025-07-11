resource "aws_iam_role" "glue_crawler_role" {
  name               = format("%s-%s-%s", var.Prefix, "glue-rds", var.EnvCode)

  
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "glue.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "glue_crawler_policy_attachment" {
  role       = aws_iam_role.glue_crawler_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSGlueServiceRole"
}

resource "aws_iam_role_policy" "glue_policy" {
  name = format("%s-%s-%s", var.Prefix, "glue", var.EnvCode)
  role = aws_iam_role.glue_crawler_role.id

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:PutLogEvents",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
             "secretsmanager:GetSecretValue",
        ]
        Effect   = "Allow"
        Resource = module.db.db_instance_master_user_secret_arn
      },
    ]
  })
}

resource "aws_glue_connection" "example" {
  name        = format("%s-%s-%s", var.Prefix, "rds", var.EnvCode)
  description = "Glue connection to RDS PostgreSQL"
  connection_type = "JDBC"

  connection_properties = {
    SECRET_ID = module.db.db_instance_master_user_secret_arn
    JDBC_CONNECTION_URL = "jdbc:postgresql://${module.db.db_instance_endpoint}:${module.db.db_instance_port}/${var.DBInstanceDatabaseName}"
  }


  # Optional: VPC configuration
  physical_connection_requirements {
    security_group_id_list = [module.db.db_subnet_group_id]
    subnet_id              =  [var.VPCDatabaseSubnetGroup]
  }
}

resource "aws_glue_catalog_database" "main" {
  name = module.db.db_instance_name
}

resource "aws_glue_crawler" "example" {
  database_name = aws_glue_catalog_database.main.name
  name          = "example"
  role          = aws_iam_role.glue_crawler_role.arn

  jdbc_target {
    connection_name = aws_glue_connection.example.name
    path            = "${var.DBInstanceDatabaseName}/public/%"
  }
}