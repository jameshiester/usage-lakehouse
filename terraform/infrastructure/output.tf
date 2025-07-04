# Output ALB DNS name for GitHub Actions job output
output "alb_dns_name" {
  value = module.ecs.alb_dns_name
}