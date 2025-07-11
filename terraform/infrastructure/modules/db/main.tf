provider "aws" {
  region = var.Region
}

data "aws_caller_identity" "current" {}
data "aws_availability_zones" "available" {}

locals {
  azs = var.AZS

  tags = {
    Environment = var.EnvTag
    EnvCode     = var.EnvCode
    Solution    = var.SolTag
  }
}

################################################################################
# RDS Module
################################################################################

module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 6"

  identifier = format("%s%s%s", var.Prefix, "rds", var.EnvCode)
  # All available versions: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html#PostgreSQL.Concepts
  engine                   = "postgres"
  engine_version           = "17.5"
  engine_lifecycle_support = "open-source-rds-extended-support-disabled"
  family                   = "postgres17" # DB parameter group
  major_engine_version     = "17"         # DB option group
  instance_class           = var.DBInstanceSize

  allocated_storage     = var.DBInstanceAllocatedStorage
  max_allocated_storage = var.DBInstanceMaxAllocatedStorage

  # NOTE: Do NOT use 'user' as the value for 'username' as it throws:
  # "Error creating DB Instance: InvalidParameterValue: MasterUsername
  # user cannot be used as it is a reserved word used by the engine"
  db_name  = var.DBInstanceDatabaseName
  username = var.DBInstanceUsername
  port     = var.DBInstancePort

  # Setting manage_master_user_password_rotation to false after it
  # has previously been set to true disables automatic rotation
  # however using an initial value of false (default) does not disable
  # automatic rotation and rotation will be handled by RDS.
  # manage_master_user_password_rotation allows users to configure
  # a non-default schedule and is not meant to disable rotation
  # when initially creating / enabling the password management feature
  manage_master_user_password_rotation              = true
  master_user_password_rotate_immediately           = false
  master_user_password_rotation_schedule_expression = "rate(30 days)"

  multi_az = true
  # db_subnet_group_name   = module.vpc.database_subnet_group
  db_subnet_group_name   = var.VPCDatabaseSubnetGroupName
  vpc_security_group_ids = [module.security_group.security_group_id]

  maintenance_window              = "Mon:00:00-Mon:03:00"
  backup_window                   = "03:00-06:00"
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
  create_cloudwatch_log_group     = true

  backup_retention_period = 1
  skip_final_snapshot     = true
  deletion_protection     = false

  performance_insights_enabled = false
  #   performance_insights_retention_period = 7
  create_monitoring_role          = true
  monitoring_interval             = 60
  monitoring_role_name            = format("%s%s%s", var.Prefix, "rds-mintoring", var.EnvCode)
  monitoring_role_use_name_prefix = true
  monitoring_role_description     = "Monitoring role for RDS database"
  apply_immediately = true

  parameters = [
    {
      name  = "autovacuum"
      value = 1
    },
    {
      name  = "rds.logical_replication"
      value = 1
      apply_method = "pending-reboot"
    },
    {
      name  = "rds.force_ssl"
      value = 0
    },
    {
      name  = "client_encoding"
      value = "utf8"
    }
  ]

  tags = local.tags
  db_option_group_tags = {
    "Sensitive" = "low"
  }
  db_parameter_group_tags = {
    "Sensitive" = "low"
  }
  cloudwatch_log_group_tags = {
    "Sensitive" = "high"
  }
}



# module "kms" {
#   source      = "terraform-aws-modules/kms/aws"
#   version     = "~> 1.0"
#   description = "KMS key for cross region automated backups replication"

#   # Aliases
#   aliases                 = [format("%s%s%s", var.Prefix, "kms", var.EnvCode)]
#   aliases_use_name_prefix = true

#   key_owners = [data.aws_caller_identity.current.arn]

#   tags = local.tags
# }

# resource "aws_db_instance_automated_backups_replication" "this" {
#   source_db_instance_arn = module.db.db_instance_arn
#   kms_key_id             = module.kms.key_arn
#   retention_period = 7
#   region = "us-west-1"
# }


################################################################################
# Supporting Resources
################################################################################



module "security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = format("%s%s%s%s", var.Prefix, "sg", var.EnvCode, "01")
  description = "PostgreSQL security group"
  vpc_id      = var.VPCID

  # ingress
  ingress_with_cidr_blocks = [
    {
      from_port   = module.db.db_instance_port
      to_port     = module.db.db_instance_port
      protocol    = "tcp"
      description = "PostgreSQL access from within VPC"
      cidr_blocks = var.VPCCIDR
    },
  ]

  tags = local.tags
}




