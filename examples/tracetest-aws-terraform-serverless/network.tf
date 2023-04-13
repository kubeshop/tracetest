module "network" {
  source                                      = "cn-terraform/networking/aws"
  name_prefix                                 = local.name
  vpc_cidr_block                              = local.vpc_cidr
  availability_zones                          = ["${local.region}a", "${local.region}b", "${local.region}c", "${local.region}d"]
  public_subnets_cidrs_per_availability_zone  = ["192.168.0.0/19", "192.168.32.0/19", "192.168.64.0/19", "192.168.96.0/19"]
  private_subnets_cidrs_per_availability_zone = ["192.168.128.0/19", "192.168.160.0/19", "192.168.192.0/19", "192.168.224.0/19"]
}

module "tracetest_alb_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = local.name
  description = "Load balancer security group"
  vpc_id      = module.network.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = 11633
      to_port     = 11633
      protocol    = "tcp"
      description = "HTTP access for Tracetest"
      cidr_blocks = "0.0.0.0/0"
      }, {
      from_port   = 16686
      to_port     = 16686
      protocol    = "tcp"
      description = "HTTP access for Jaeger UI"
      cidr_blocks = "0.0.0.0/0"
      }, {
      from_port   = 16685
      to_port     = 16685
      protocol    = "tcp"
      description = "HTTP access for Jaeger API"
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

resource "aws_lb" "tracetest-alb" {
  name               = "tracetest-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [module.tracetest_alb_security_group.security_group_id]
  subnets            = module.network.public_subnets_ids

  enable_deletion_protection = false
  tags                       = local.tags
}

// INTERNAL ALB
module "internal_tracetest_alb_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = local.name
  description = "Internal Load balancer security group"
  vpc_id      = module.network.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = 16685
      to_port     = 16685
      protocol    = "tcp"
      description = "HTTP access for Jaeger API"
      cidr_blocks = local.vpc_cidr
      }, {
      from_port   = 4318
      to_port     = 4318
      protocol    = "tcp"
      description = "HTTP access for Jaeger Collector"
      cidr_blocks = local.vpc_cidr
  }]

  egress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 65535
      protocol    = "-1"
      description = "HTTP access to Anywhere"
      cidr_blocks = "0.0.0.0/0"
  }]
}

resource "aws_lb" "internal_tracetest_alb" {
  name               = "tracetest-internal-alb"
  internal           = true
  load_balancer_type = "application"
  security_groups    = [module.internal_tracetest_alb_security_group.security_group_id]
  subnets            = module.network.private_subnets_ids

  enable_deletion_protection = false
  tags                       = local.tags
}

module "lambda_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = "${local.name}_lambda_security_group"
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

module "tracetest_ecs_service_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = "tracetest_ecs_service_security_group"
  description = "ECS Service security group"
  vpc_id      = module.network.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 65535
      protocol    = "tcp"
      description = "HTTP access from VPC"
      cidr_blocks = local.vpc_cidr
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