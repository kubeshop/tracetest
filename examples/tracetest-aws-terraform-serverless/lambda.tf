resource "aws_lambda_function" "hello_world" {
  function_name = "HelloWorld"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_hello_world.key

  runtime = "nodejs12.x"
  handler = "hello.handler"
  timeout = 30

  environment {
    variables = {
      NODE_OPTIONS    = "--require tracing.js"
      OTEL_EXPORTER_OTLP_ENDPOINT = "http://${aws_lb.internal_tracetest_alb.dns_name}:4318"
    }
  }

  source_code_hash = data.archive_file.lambda_hello_world.output_base64sha256

  vpc_config {
    subnet_ids         = module.network.public_subnets_ids
    security_group_ids = [module.lambda_security_group.security_group_id]
  }

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "hello_world" {
  name = "/aws/lambda/${aws_lambda_function.hello_world.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "lambda_exec" {
  name = "serverless_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_network_policy_attachment" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_network_policy.arn
}

resource "aws_iam_policy" "lambda_network_policy" {
  name   = "lambda_network_policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeNetworkInterfaces",
        "ec2:CreateNetworkInterface",
        "ec2:DeleteNetworkInterface",
        "ec2:DescribeInstances",
        "ec2:AttachNetworkInterface"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}


module "lambda_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = "lambda_security_group"
  description = "Lambda security group"
  vpc_id      = module.network.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 65535
      protocol    = "-1"
      description = "HTTP access from anywhere"
      cidr_blocks = "0.0.0.0/0"
  }]

  egress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 65535
      protocol    = "-1"
      description = "HTTP access to anywhere"
      cidr_blocks = "0.0.0.0/0"
  }]
}
