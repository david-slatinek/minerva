output "db_endpoint" {
  description = "The connection endpoint in address:port format."
  value       = aws_db_instance.database.address
}
