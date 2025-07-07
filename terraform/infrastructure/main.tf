terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
  backend "s3" {}
}



provider "aws" {
  region = var.Region
}

data "aws_availability_zones" "available" {}

locals {
  vpc_cidr                     = "10.0.0.0/16"
  azs                          = slice(data.aws_availability_zones.available.names, 0, 3)
  preferred_maintenance_window = "sun:05:00-sun:06:00"

  tags = {
    EnvCode     = var.EnvCode
    Environment = var.EnvTag
    Solution    = var.SolTag
  }
}


module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.0"

  name               = format("%s%s%s%s", var.Prefix, "vpc", var.EnvCode, "01")
  cidr               = local.vpc_cidr
  enable_nat_gateway = true
  single_nat_gateway = true

  enable_dns_hostnames = true
  enable_dns_support   = true

  azs              = local.azs
  public_subnets   = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k)]
  private_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 3)]
  database_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 6)]
  intra_subnets    = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 9)]

  create_database_subnet_group = true

  tags = local.tags
}

resource "aws_security_group" "lambda" {
  name        = format("%s%s%s%s", var.Prefix, "scg", var.EnvCode, "lambda")
  description = " Application Security Group"
  vpc_id      = module.vpc.vpc_id

  tags = {
    Name         = format("%s%s%s-%s", var.Prefix, "scg", var.EnvCode, "lambda")
    resourcetype = "security"
    codeblock    = "network-3tier"
  }
}

resource "aws_vpc_security_group_egress_rule" "postgres" {
  security_group_id = aws_security_group.lambda.id
  description       = "postgres"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 5432
  ip_protocol       = "tcp"
  to_port           = 5432
}

resource "aws_vpc_security_group_egress_rule" "https" {
  security_group_id = aws_security_group.lambda.id
  description       = "https"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 443
  ip_protocol       = "tcp"
  to_port           = 443
}

module "vpc_endpoints" {
  source = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"

  vpc_id = module.vpc.vpc_id

  create_security_group      = true
  security_group_name_prefix = format("%s%s-", "vpc-endpoints-", var.EnvCode)
  security_group_description = "VPC endpoint security group"
  security_group_rules = {
    ingress_https = {
      description = "HTTPS from VPC"
      cidr_blocks = [module.vpc.vpc_cidr_block]
    }
  }

  endpoints = {
    ssm = {
      service             = "ssm"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
    },
    ecs = {
      service             = "ecs"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      subnet_configurations = [
        for v in module.vpc.private_subnet_objects :
        {
          ipv4      = cidrhost(v.cidr_block, 10)
          subnet_id = v.id
        }
      ]
    },
    ecs_telemetry = {
      create              = false
      service             = "ecs-telemetry"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
    },
  }

  tags = merge(local.tags, {
    Project  = "Secret"
    Endpoint = "true"
  })
}

module "db" {
  source                 = "./modules/db"
  Region                 = var.Region
  Prefix                 = var.Prefix
  VPCID                  = module.vpc.vpc_id
  VPCDatabaseSubnetGroup = module.vpc.database_subnet_group
  GitHubRepo             = var.GitHubRepo
  SolTag                 = var.SolTag
  DBInstanceSize         = var.DBInstanceSize
  EnvCode                = var.EnvCode
  VPCCIDR                = local.vpc_cidr
  AZS                    = local.azs
  EnvTag                 = var.EnvTag
}

module "ecs" {
  source             = "./modules/ecs"
  Region             = var.Region
  Prefix             = var.Prefix
  VPCID              = module.vpc.vpc_id
  SolTag             = var.SolTag
  EnvCode            = var.EnvCode
  VPCCIDR            = local.vpc_cidr
  AZS                = local.azs
  EnvTag             = var.EnvTag
  DBSecretArn        = module.db.db_instance_master_user_secret_arn
  DBHost             = module.db.db_instance_address
  DBInstancePort     = module.db.db_instance_port
  DBInstanceUsername = module.db.db_instance_username
  ECRRepo            = var.ECRRepo
  PublicSubnets      = module.vpc.public_subnets
  PrivateSubnets     = module.vpc.private_subnets
  ECSCluster         = var.ECSCluster
  ECSService         = var.ECSService
}