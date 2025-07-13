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

variable "PublicSubnets" {
  description = "IDs of public subnets"
  type        = list(string)
}

variable "PrivateSubnets" {
  description = "IDs of private subnets"
  type        = list(string)
}

# Networking
variable "VPCDatabaseSubnetGroup" {
  description = "Database subnet group id created in the VPC"
  type        = string
}

# Networking
variable "VPCDatabaseSubnetGroupName" {
  description = "Database subnet group name created in the VPC"
  type        = string
}

# Tagging and naming
variable "Prefix" {
  description = "Prefix used to name all resources"
  type        = string
}

variable "DBInstanceSize" {
  description = "Size of db instance to create"
  type        = string
}

variable "DBInstancePort" {
  description = "Port the instance will use to communicate"
  type        = string
  default     = 5432
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
}

variable "DBInstanceAllocatedStorage" {
  description = "Allocated storage for the database in GB"
  type        = number
  default     = 20
}

variable "DBInstanceMaxAllocatedStorage" {
  description = "Maximum Allocated storage for the database in GB"
  type        = number
  default     = 100
}

variable "AZS" {
  description = "A list of availability zones names or ids in the region"
  type        = list(string)
}

variable "SolTag" {
  description = "Solution tag value. All resources are created with a 'Solution' tag name and the value you set here"
  type        = string
}
variable "GitHubRepo" {
  description = "GitHub repository name"
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