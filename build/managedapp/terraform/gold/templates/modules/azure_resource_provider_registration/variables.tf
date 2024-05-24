variable "platform" {
  type        = string
  description = <<-EOL
  The current OS, determined by Terragrunt. Some of the returned values can be:
  darwin, freebsd, linux, windows
  See: https://terragrunt.gruntwork.io/docs/reference/built-in-functions/#get_platform
  EOL
}

variable "name" {
  type        = string
  description = "Name of resource provider to register"
}

variable "subscription_id" {
  type        = string
  description = "ID for the Azure subscription to register resource providers in"
}
