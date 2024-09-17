resource "aws_ecr_repository" "reg" {
  name                 = var.name
  image_tag_mutability = "MUTABLE"
}
