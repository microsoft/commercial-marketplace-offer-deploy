provider "azurerm" {
  features {}
}

variable "location" {
  type = string
}

variable "tier" {
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
