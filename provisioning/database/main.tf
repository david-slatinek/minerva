locals {
  db_name      = "minerva"
  project_name = "${local.db_name}-project"
}

resource "aws_db_instance" "database" {
  allocated_storage           = 20
  allow_major_version_upgrade = false

  db_name        = local.db_name
  engine         = "postgres"
  engine_version = "16.4"
  identifier     = local.project_name
  instance_class = "db.t4g.micro"

  username = var.username
  password = var.password

  publicly_accessible = true
  skip_final_snapshot = true

  tags = {
    Name        = local.project_name
    Environment = var.env
  }
}
