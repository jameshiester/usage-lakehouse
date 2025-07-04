# Regions
variable "Region" {
  description = "AWS depoloyment region"
  type        = string
}

# Networking
variable "VPCCIDR" {
  description = "VPC CIDR range"
  type        = string
}

# Networking
variable "VPCID" {
  description = "ID of the VPC"
  type        = string
}

# Tagging and naming
variable "Prefix" {
  description = "Prefix used to name all resources"
  type        = string
}

variable "AZS" {
  description = "A list of availability zones names or ids in the region"
  type        = list(string)
}

variable "PublicSubnets" {
  description = "A list of public subnets"
  type        = list(string)
}

variable "PrivateSubnets" {
  description = "A list of private subnets"
  type        = list(string)
}

variable "SolTag" {
  description = "Solution tag value. All resources are created with a 'Solution' tag name and the value you set here"
  type        = string
}

variable "EnvCode" {
  description = "2 character code used to name all resources e.g. 'pd' for production"
  type        = string
}

variable "EnvTag" {
  description = "Environment tag value. All resources are created with an 'Environment' tag name and the value you set here"
  type        = string
}

# Web App Build
variable "ECRRepo" {
  description = "Name of Amazon ECR repository"
  type        = string
}
variable "ImageTag" {
  description = "Amazon ECR Microsoft sample application Image Tag"
  type        = string
}

variable "DBSecretArn" {
  description = "ARN of the secret used to hold the RDS password"
  type        = string
}

variable "DBInstancePort" {
  description = "Port the instance will use to communicate"
  type        = string
}

variable "DBInstanceDatabaseName" {
  description = "Name of the default database"
  type        = string
  default     = "postgres"
}

variable "DBInstanceUsername" {
  description = "Name of the default user"
  type        = string
  default     = "postgres"
  sensitive   = true
}

variable "DBHost" {
  description = "Username of the DB user"
  type        = string
}