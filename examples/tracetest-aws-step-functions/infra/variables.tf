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
  name            = "tracetest-step-functions"
  region          = var.aws_region
  tracetest_image = "kubeshop/tracetest:${var.tracetest_version}"
  environment     = var.environment

  db_name     = "postgres"
  db_username = "postgres"

  vpc_cidr = "192.168.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  tags = {
    Name    = local.name
    Example = local.name
  }
}
