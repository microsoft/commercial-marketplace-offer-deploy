provider "azurerm" {
  features {}
}

module "storage" {
  source = "./child-terraform"

  location = "East US"
  resource_group_name = "bobjacterraform"
}
