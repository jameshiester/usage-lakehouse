variable "GitHubRepo" {
  description = "GitHub repository name"
  type        = string
  default     = "usage-lakehouse"
}

variable "GitHubOrg" {
  description = "GitHub Organization Name / User Name"
  type        = string
  default     = "jameshiester"
}

# Regions
variable "Region" {
  description = "AWS deployment region"
  type        = string
  default     = "us-east-1"
}

variable "Prefix" {
  description = "Prefix used to name all resources"
  type        = string
  default     = "uwh"
}