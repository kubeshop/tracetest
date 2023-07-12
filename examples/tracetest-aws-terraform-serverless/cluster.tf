resource "aws_ecs_cluster" "tracetest-cluster" {
  name = "tracetest"
  tags = local.tags
}

resource "aws_iam_role" "tracetest_task_execution_role" {
  name = "tracetest_task_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = ["ecs-tasks.amazonaws.com", "ecs.amazonaws.com"]
        }
      },
    ]
  })
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceRole"]

  tags = local.tags
}

resource "aws_iam_role_policy" "tracetest_task_execution_role_policy" {
  name = "tracetest_task_execution_role_policy"
  role = aws_iam_role.tracetest_task_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:PutLogEvents",
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:DescribeLogStreams",
          "logs:DescribeLogGroups",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}
