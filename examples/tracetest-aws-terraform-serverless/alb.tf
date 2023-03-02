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
