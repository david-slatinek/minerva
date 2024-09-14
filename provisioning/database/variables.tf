variable "access_key" {
  type        = string
  description = "AWS Access Key"
  sensitive   = true
}

variable "secret_key" {
  type        = string
  description = "AWS Secret Key"
  sensitive   = true
}

variable "password" {
  type        = string
  description = "Root password"
  sensitive   = true
}

variable "username" {
  type        = string
  description = "Root username"
  sensitive   = true
}

variable "env" {
  type        = string
  description = "The environment name"
  default     = "prod"

  validation {
    condition     = can(regex("dev|test|prod", var.env))
    error_message = "Env must be one of the following: 'dev', 'test' or 'prod'."
  }
}
