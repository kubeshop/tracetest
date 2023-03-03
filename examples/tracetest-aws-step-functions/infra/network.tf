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

  ingress_with_cidr_blocks = [{
    from_port   = 11633
    to_port     = 11633
    protocol    = "tcp"
    description = "HTTP access for Tracetest"
    cidr_blocks = "0.0.0.0/0"
  }]

  egress_with_cidr_blocks = [{
    from_port   = 0
    to_port     = 65535
    protocol    = "-1"
    description = "HTTP access to anywhere"
    cidr_blocks = "0.0.0.0/0"
  }]
}

resource "aws_lb" "tracetest_alb" {
  name               = "${local.name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [module.tracetest_alb_security_group.security_group_id]
  subnets            = module.network.public_subnets_ids

  enable_deletion_protection = false
  tags                       = local.tags
}

