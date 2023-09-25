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


resource "azurerm_virtual_network" "vNet" {
  name                = var.Vnet_Name
  address_space       = var.Vnet_AddressSpace
  location            = var.location
  resource_group_name = var.rgName
}

resource "azurerm_network_security_group" "NSG" {
  name                = var.NSG_Name
  location            = var.location
  resource_group_name = var.rgName

  security_rule {
    name                       = "allowRDP"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3389"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
 

resource "azurerm_subnet" "default" {
  name                 = "default"
  resource_group_name  = var.rgName
  virtual_network_name = azurerm_virtual_network.vNet.name
  address_prefixes     = var.SubnetPrefix
}


resource "azurerm_subnet_network_security_group_association" "associateNSG" {
  subnet_id                 = azurerm_subnet.default.id
  network_security_group_id = azurerm_network_security_group.NSG.id
}


#build Pub IP resources for East 
resource "azurerm_public_ip" "PubIP" {
  name                = var.public_IP
  resource_group_name = var.rgName
  location            = var.location
  allocation_method   = "Dynamic"  
}


resource "azurerm_network_interface" interface_name {
  name                = var.interface_name
  location            = var.location
  resource_group_name = var.rgName
 
  ip_configuration {
    name                          = "interfaceconfiguration"
    subnet_id                     = azurerm_subnet.default.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id = azurerm_public_ip.PubIP.id
  }
}