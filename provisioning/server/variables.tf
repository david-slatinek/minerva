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

variable "filepath" {
  type        = string
  description = "Path to the public key"
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

variable "elastic_ip" {
  type        = string
  description = "Already created elastic ip"
}
