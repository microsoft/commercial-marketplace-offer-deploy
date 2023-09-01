resource "azurerm_resource_group" "example" {
  name     = "modm-dev-vmi"
  location = "East US"
}

resource "azurerm_storage_account" "example" {
  name                     = "bobjacterra2"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
