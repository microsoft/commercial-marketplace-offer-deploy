provider "azurerm" {
  features {}
}

variable "resource_group_name" {
  description = "The name of the resource group in which all resources will be deployed"
  type        = string
}

variable "location" {
  description = "Azure region for the resources"
  default     = "West US"
}

resource "azurerm_virtual_network" "example_vnet" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = var.location
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet" "example_subnet" {
  name                 = "example-subnet"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.example_vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_interface" "example_nic" {
  name                = "example-nic"
  location            = var.location
  resource_group_name = var.resource_group_name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example_subnet.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "example_vm" {
  name                  = "example-vm"
  location              = var.location
  resource_group_name   = var.resource_group_name
  network_interface_ids = [azurerm_network_interface.example_nic.id]
  vm_size               = "Standard_F2"

  storage_os_disk {
    name              = "osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  os_profile {
    computer_name  = "examplevm"
    admin_username = "adminuser"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "test"
  }
}

resource "azurerm_storage_account" "example_sa" {
  name                     = "examplestoracc"
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "testing"
  }
}

resource "azurerm_public_ip" "example_pip" {
  name                = "example-pip"
  location            = var.location
  resource_group_name = var.resource_group_name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "example_nsg" {
  name                = "example-nsg"
  location            = var.location
  resource_group_name = var.resource_group_name
}

resource "azurerm_network_security_rule" "example_nsr" {
  name                        = "SSH"
  priority                    = 1001
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = var.resource_group_name
  network_security_group_name = azurerm_network_security_group.example_nsg.name
}

resource "azurerm_app_service_plan" "example_asp" {
  name                = "example-asp"
  location            = var.location
  resource_group_name = var.resource_group_name

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "example_app_service" {
  name                = "example-app-service"
  location            = var.location
  resource_group_name = var.resource_group_name
  app_service_plan_id = azurerm_app_service_plan.example_asp.id

  site_config {
    dotnet_framework_version = "v5.0"
  }

  app_settings = {
    "SOME_KEY" = "some-value"
  }
}

resource "azurerm_sql_server" "example_sql_server" {
  name                         = "examplesqlserver"
  resource_group_name          = var.resource_group_name
  location                     = var.location
  version                      = "12.0"
  administrator_login          = "adminuser"
  administrator_login_password = "Password1234!"

  tags = {
    environment = "testing"
  }
}

resource "azurerm_sql_database" "example_sql_db" {
  name                = "examplesqldb"
  resource_group_name = var.resource_group_name
  location            = var.location
  server_name         = azurerm_sql_server.example_sql_server.name

  tags = {
    environment = "testing"
  }
}

resource "azurerm_cosmosdb_account" "example_cosmosdb" {
  name                = "example-cosmosdb"
  location            = var.location
  resource_group_name = var.resource_group_name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_automatic_failover = false

  capabilities {
    name = "EnableCassandra"
  }

  consistency_policy {
    consistency_level       = "Session"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = var.location
    failover_priority = 0
  }
}