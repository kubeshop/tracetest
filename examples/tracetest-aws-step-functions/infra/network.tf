module "network" {
  source                                      = "cn-terraform/networking/aws"
  name_prefix                                 = local.name
  vpc_cidr_block                              = local.vpc_cidr
  availability_zones                          = ["${local.region}a", "${local.region}b", "${local.region}c", "${local.region}d"]
  public_subnets_cidrs_per_availability_zone  = ["192.168.0.0/19", "192.168.32.0/19", "192.168.64.0/19", "192.168.96.0/19"]
  private_subnets_cidrs_per_availability_zone = ["192.168.128.0/19", "192.168.160.0/19", "192.168.192.0/19", "192.168.224.0/19"]
}
