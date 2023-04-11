terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.55.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "4.0.4"
    }
  }
}

provider "aws" {
  region = local.region
}

data "aws_caller_identity" "current" {}
data "aws_availability_zones" "available" {}

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

  provisioning = <<EOF
---
type: PollingProfile
spec:
  name: default
  strategy: periodic
  default: true
  periodic:
    retryDelay: 5s
    timeout: 10m

---
type: DataStore
spec:
  name: jaeger
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
