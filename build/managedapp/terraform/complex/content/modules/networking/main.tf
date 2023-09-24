locals {
    timestamp_suffix  = formatdate("YYYYMMDDHHmmss", timestamp())
    vnet_name         = "modmvnet-${local.timestamp_suffix}"
    subnet_name       = "modmsubnet-${local.timestamp_suffix}"
}

resource "azurerm_virtual_network" "vnet" {
  name                = local.vnet_name
  address_space       = ["10.0.0.0/16"]
  location            = var.location
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet" "subnet" {
  name                 = local.subnet_name
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

