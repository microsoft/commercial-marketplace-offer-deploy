resource "azurerm_storage_account" "storage_acc" {
  name                     = "stg${var.storage_suffix}"
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

