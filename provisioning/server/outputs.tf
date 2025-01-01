output "server_ip" {
  description = "Public IP address of the created server instance"
  value       = aws_instance.api.public_ip
}

output "server_dns" {
  description = "Public DNS name assigned to the server instance"
  value       = aws_instance.api.public_dns
}

output "elastic_ip" {
  description = "Elastic IP address of the created server instance"
  value       = aws_eip_association.api_eip_association.public_ip
}
