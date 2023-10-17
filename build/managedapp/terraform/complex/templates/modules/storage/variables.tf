variable "location" {
  description = "Azure region for the resources"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group in which resources will be deployed"
  type        = string
}

variable "storage_suffix" {
  description = "Suffix to ensure storage account name uniqueness"
  type        = string
}
