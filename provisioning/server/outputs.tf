output "server_ip" {
  description = "Public IP address of created server instance"
  value       = aws_instance.api.public_ip
}

output "server_dns" {
  description = "Public DNS name assigned to the server instance"
  value       = aws_instance.api.public_dns
}

output "elastic_ip" {
  description = "Elastic IP address of created server instance"
  value       = aws_eip.eip.public_ip
}

output "elastic_dns" {
  description = "Elastic DNS name assigned to server instance"
  value       = aws_eip.eip.public_dns
}
