provider "azurerm" {
  features {}
}

variable "location" {
  type = string
}

variable "resource_group_name" {
  type = string
}

module "storage" {
  source = "./storage"

  location = var.location
  resource_group_name = var.resource_group_name
}
