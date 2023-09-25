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

resource "azurerm_windows_virtual_machine" "vm" {
  name                  = var.computerName
  location              = var.location
  resource_group_name   = var.rgName
  admin_username = var.userAccount
  admin_password = var.vmPasswrd
  network_interface_ids = var.interface_name 
  size                  = "Standard_DS1_v2"
 

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  os_disk {
    name              = "osdisk1"
    caching           = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
}