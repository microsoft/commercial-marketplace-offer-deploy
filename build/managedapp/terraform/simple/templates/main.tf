provider "azurerm" {
  features {}
  skip_provider_registration=true
}

variable "artifactsLocationSasToken" {
  type = string
}

variable "location" {
  type = string
}

variable "resourceGroupName" {
  type = string
}

module "storage" {
  source = "./storage"

  location = var.location
  resource_group_name = var.resourceGroupName
}
