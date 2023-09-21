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
  subscription_id = var.subscriptionID
  tenant_id       = var.tenantID
  features {
  }
}

resource "azurerm_virtual_network_peering" "peer" {
  name                      = var.peerName
  resource_group_name = var.rgName
  virtual_network_name =  var.localPeerVnetName
  remote_virtual_network_id = var.remotePeerVnetID
}