locals {
  azs = var.AZS

  tags = {
    Environment = var.EnvTag
    EnvCode     = var.EnvCode
    Solution    = var.SolTag
  }
}

resource "aws_ecs_task_definition" "mswebapp" {
  family                   = format("%s%s%s", var.Prefix, "ect", var.EnvCode)
  requires_compatibilities = ["FARGATE"]
  cpu                      = 1024
  memory                   = 2048
  network_mode             = "awsvpc"
  track_latest             = true
  execution_role_arn       = aws_iam_role.ecstaskexec.arn
  task_role_arn            = aws_iam_role.ecstask.arn
  container_definitions = jsonencode([
    {
      name                   = "mswebapp"
      image                  = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.Region}.amazonaws.com/${var.ECRRepo}:${var.ImageTag}"
      cpu                    = 256
      memory                 = 512
      essential              = true
      readonlyRootFilesystem = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]
      environment = [{
        name  = "POSTGRES_DB"
        value = var.DBInstanceDatabaseName
        }, {
        name  = "POSTGRES_HOST"
        value = var.DBHost
        }, {
        name  = "DB_MASTER_SECRET_ARN"
        value = var.DBSecretArn
        },
         {
        name  = "POSTGRES_USER"
        value = var.DBInstanceUsername
        }
        , {
          name  = "POSTGRES_PORT"
          value = var.DBInstancePort
      }]
      logconfiguration = {
        logDriver = "awslogs",
        options = {
          awslogs-group         = "${aws_cloudwatch_log_group.mswebapp.name}",
          awslogs-region        = "${var.Region}",
          awslogs-stream-prefix = "awslogs-"
        }
      }
      healthCheck = {
        command         = ["CMD-SHELL", "curl -f http://localhost:8080/healthz || exit 1"]
        intervalSeconds = 30
        timeoutSeconds  = 5
        retries         = 3
        startPeriod     = 30
      }
    }
  ])
}


# Create Amazon ECS task service
resource "aws_ecs_service" "mswebapp" {
  name            = format("%s%s%s%s", var.Region, "iar", var.EnvCode, "api")
  cluster         = aws_ecs_cluster.mswebapp.id
  task_definition = aws_ecs_task_definition.mswebapp.arn
  launch_type     = "FARGATE"
  desired_count   = 2
  propagate_tags  = "TASK_DEFINITION"


  network_configuration {
    subnets         = var.PrivateSubnets
    security_groups = [aws_security_group.app01.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.mswebapp.arn
    container_name   = "mswebapp"
    container_port   = 8080
  }

  tags = {
    Name  = format("%s%s%s%s", var.Region, "iar", var.EnvCode, "api")
    rtype = "ecsservice"
  }
}