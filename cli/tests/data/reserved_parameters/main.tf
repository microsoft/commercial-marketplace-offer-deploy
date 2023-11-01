provider "azurerm" {
  features {}
}

variable "location" {
  type = string
}

// this is the reserved parameter that installer provided
variable "resourceGroupName" {
  type = string
}

module "storage" {
  source = "./storage"

  location = var.location
  resource_group_name = var.resourceGroupName
}
