terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "=2.99.0"
    }
  }
}

provider "azurerm" {
  # Configuration options
  subscription_id = var.subscription_ID
  tenant_id       = var.tenant_ID
  features {
  }
}


# Create a resource group #
resource "azurerm_resource_group" "rg" {
  name     = var.rg_Name
  location = var.rg_Location
}




