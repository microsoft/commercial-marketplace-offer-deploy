provider "azurerm" {
  features {}
  skip_provider_registration = true
}

variable "resourceGroupName" {
  description = "The name of the resource group in which all resources will be deployed"
  type        = string
}

variable "location" {
  description = "Azure region for the resources"
  type = string
  default     = "West US"
}

variable "sql_admin_password" {
  description = "The password for the SQL administrator login"
  type        = string
  sensitive   = true
}

variable "sql_admin_username" {
  description = "The username for the SQL administrator login"
  type        = string
}


locals {
  timestamp_suffix  = formatdate("YYYYMMDDHHmmss", timestamp())
  storage_name      = "stg${substr(local.timestamp_suffix, 8, 6)}" # Will look like stg123456 (Azure Storage Account names must be between 3 and 24 characters in length, use lowercase letters and numbers only)
  sql_server_name   = "sqlsrv-${local.timestamp_suffix}"
  sql_db_name       = "sqldb-${local.timestamp_suffix}"
  cosmosdb_name     = "cosmosdb-${local.timestamp_suffix}"
  app_service_name  = "appsvc-${local.timestamp_suffix}"
  app_service_plan_name  = "appsvcplan-${local.timestamp_suffix}"
  vnet_name         = "modmvnet-${local.timestamp_suffix}"
  subnet_name       = "modmsubnet-${local.timestamp_suffix}"
  nic_name          = "modmnic-${local.timestamp_suffix}"
  vm_name           = "modmvm-${local.timestamp_suffix}"
  pip_name          = "modmpip-${local.timestamp_suffix}"
  nsg_name          = "modmnsg-${local.timestamp_suffix}"
  storage_name_suffix  = formatdate("YYYYMMDDHHmmss", timestamp())
}

module "networking" {
  source             = "./modules/networking"
  location           = var.location
  resource_group_name = var.resourceGroupName
}

module "storage" {
  source             = "./modules/storage"
  location           = var.location
  resource_group_name = var.resourceGroupName
  storage_suffix     = local.storage_name_suffix
}

resource "azurerm_network_interface" "example_nic" {
  name                = local.nic_name
  location            = var.location
  resource_group_name = var.resourceGroupName

  ip_configuration {
    name                          = "internal"
    subnet_id                     = module.networking.subnet_id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "example_vm" {
  name                  = local.vm_name
  location              = var.location
  resource_group_name   = var.resourceGroupName
  network_interface_ids = [azurerm_network_interface.example_nic.id]
  vm_size               = "Standard_F2"

  storage_os_disk {
    name              = "osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
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
  name                     = local.storage_name
  resource_group_name      = var.resourceGroupName
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "testing"
  }
}

resource "azurerm_public_ip" "example_pip" {
  name                = local.pip_name
  location            = var.location
  resource_group_name = var.resourceGroupName
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "example_nsg" {
  name                = local.nsg_name
  location            = var.location
  resource_group_name = var.resourceGroupName
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
  resource_group_name         = var.resourceGroupName
  network_security_group_name = azurerm_network_security_group.example_nsg.name
}

resource "azurerm_service_plan" "example_asp" {
  name                = local.app_service_plan_name
  location            = var.location
  resource_group_name = var.resourceGroupName
  os_type = "Linux"
  sku_name = "S1"
}

resource "azurerm_app_service" "example_app_service" {
  name                = local.app_service_name
  location            = var.location
  resource_group_name = var.resourceGroupName
  app_service_plan_id = azurerm_service_plan.example_asp.id

  site_config {
    dotnet_framework_version = "v5.0"
  }

  app_settings = {
    "SOME_KEY" = "some-value"
  }
}

resource "azurerm_mssql_server" "example_sql_server" {
  name                         = local.sql_server_name
  resource_group_name          = var.resourceGroupName
  location                     = var.location
  version                      = "12.0"
  administrator_login          = "adminuser"
  administrator_login_password = var.sql_admin_password

  tags = {
    environment = "testing"
  }
}

resource "azurerm_sql_database" "example_sql_db" {
  name                = local.sql_db_name
  resource_group_name = var.resourceGroupName
  location            = var.location
  server_name         = azurerm_mssql_server.example_sql_server.name

  tags = {
    environment = "testing"
  }
}

resource "azurerm_cosmosdb_account" "example_cosmosdb" {
  name                = local.cosmosdb_name
  location            = var.location
  resource_group_name = var.resourceGroupName
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