provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_resource_group" "rg" {
  name = "bobjacterraform"
}

data "azurerm_resource_group" "example" {
  name = "bobjacterraform2"
}

module "storage" {
  source               = "./child-terraform"
  location             = "East US"
  resource_group_name  = data.azurerm_resource_group.rg.name
}
