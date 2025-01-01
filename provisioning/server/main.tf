locals {
  name     = "minerva"
  key_name = "${local.name}-key"
  sg       = "${local.name}-sg"
}

data "aws_ami" "debian" {
  most_recent = true
  owners      = ["136693071363"]

  filter {
    name   = "name"
    values = ["debian-12-amd64-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
}

resource "aws_key_pair" "key_pair" {
  key_name   = local.key_name
  public_key = file(var.filepath)

  tags = {
    Name        = local.key_name
    Environment = var.env
  }
}

resource "aws_security_group" "security_group" {
  name = local.sg

  tags = {
    Name        = local.sg
    Environment = var.env
  }
}

data "http" "my_ip" {
  url = "https://ipv4.icanhazip.com"
}


resource "aws_vpc_security_group_ingress_rule" "allow_ssh" {
  security_group_id = aws_security_group.security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  ip_protocol       = "tcp"
  to_port           = 22

  tags = {
    Name        = "${local.name}-ssh"
    Protocol    = "ssh"
    Environment = var.env
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_http" {
  security_group_id = aws_security_group.security_group.id
  cidr_ipv4         = "${chomp(data.http.my_ip.response_body)}/32"
  from_port         = 80
  ip_protocol       = "tcp"
  to_port           = 80

  tags = {
    Name        = "${local.name}-http"
    Protocol    = "http"
    Environment = var.env
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_api_postgres" {
  security_group_id = aws_security_group.security_group.id
  cidr_ipv4         = "${chomp(data.http.my_ip.response_body)}/32"
  from_port         = 5432
  ip_protocol       = "tcp"
  to_port           = 5432

  tags = {
    Name        = "${local.name}-api-postgres"
    Protocol    = "postgres"
    Environment = var.env
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_cli_postgres" {
  security_group_id = aws_security_group.security_group.id
  cidr_ipv4         = "${chomp(data.http.my_ip.response_body)}/32"
  from_port         = 5000
  ip_protocol       = "tcp"
  to_port           = 5000

  tags = {
    Name        = "${local.name}-cli-postgres"
    Protocol    = "postgres"
    Environment = var.env
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_netdata" {
  security_group_id = aws_security_group.security_group.id
  cidr_ipv4         = "${chomp(data.http.my_ip.response_body)}/32"
  from_port         = 19999
  ip_protocol       = "tcp"
  to_port           = 19999

  tags = {
    Name        = "${local.name}-netdata"
    Protocol    = "netdata"
    Environment = var.env
  }
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4" {
  security_group_id = aws_security_group.security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"

  tags = {
    Name        = "${local.name}-all"
    Protocol    = "all"
    Environment = var.env
  }
}

resource "aws_instance" "api" {
  ami                         = data.aws_ami.debian.id
  instance_type               = "t3.medium"
  associate_public_ip_address = true
  key_name                    = aws_key_pair.key_pair.key_name

  vpc_security_group_ids = [aws_security_group.security_group.id]

  root_block_device {
    volume_size = 30
    tags = {
      Name        = "${local.name}-storage"
      Environment = var.env
    }
  }

  tags = {
    Name        = local.name
    Environment = var.env
  }
}

data "aws_eip" "api_eip" {
  public_ip = var.elastic_ip

  tags = {
    Name = "${local.name}-eip"
  }
}

resource "aws_eip_association" "api_eip_association" {
  instance_id   = aws_instance.api.id
  allocation_id = data.aws_eip.api_eip.id
}
