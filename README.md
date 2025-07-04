# Usage Lakehouse

## Purpose
The purpose of this repository is to help me experiment with AWS's Big Data and AI offerings.  In particular I wanted to test the performance & costs of Apache Iceberg vs. a traditional RDS (Postgres) database.  I also wanted to gain experience with AI development tools including Cursor.  

## Setup

Note: These directions have been summarized as the purpose of this repository is more for self-exploration and is not intended for other developers.  I beleive a thoughtful README is important when working on a team project to assist with onboarding and documenting different processes.

1. Create an AWS account and configure a default profile using the CLI.
2. Create an S3 bucket for the default terraform state backend.
3. Create a dynamodb table for managing the terraform lock
2. Install terraform
3. Install Go
4. Run `cd ./terraform/repository`.  This directory includes the IAM role and configuration that will allow Github Actions to create the actual infrastructure needed to run the application.  It also creates branches, branch protections, and environment variables for running the CI/CD pipelines.
5. Run `terraform apply --auto-approve`