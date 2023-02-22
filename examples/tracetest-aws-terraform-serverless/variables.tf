variable "aws_region" {
  type    = string
  default = "us-west-2"
}

variable "tracetest_version" {
  type    = string
  default = "latest"
}

variable "environment" {
  type    = string
  default = "dev"
}

locals {
  name            = "tracetest"
  region          = var.aws_region
  tracetest_image = "kubeshop/tracetest:${var.tracetest_version}"
  environment     = var.environment

  db_name     = "postgres"
  db_username = "postgres"

  vpc_cidr = "192.168.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  privisioning = <<EOF
dataStore:
  type: jaeger
  jaeger:
    endpoint: ${aws_lb.internal_tracetest_alb.dns_name}:16685
    tls:
      insecure_skip_verify: true
  EOF

  tags = {
    Name    = local.name
    Example = local.name
  }
}
