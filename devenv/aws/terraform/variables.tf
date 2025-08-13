variable "key_pair_name" {
  description = "The name of the SSH key pair to use for EC2 instances"
  type        = string
}

variable "ec2_profile" {
  description = "The AWS profile to use for EC2 instances"
  type        = string
}

variable "datadog_api_key" {
  description = "Datadog API Key"
  type        = string
}

variable "datadog_site" {
  description = "Datadog site (e.g., datadoghq.com)"
  type        = string
  
}