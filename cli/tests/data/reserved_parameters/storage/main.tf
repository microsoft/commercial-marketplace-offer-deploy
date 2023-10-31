variable "location" {
  type = string
}

variable "resource_group_name" {
  type = string
}

resource "random_id" "storage_account_name_unique" {
  byte_length = 8
}

resource "azurerm_storage_account" "example" {
  name                     = random_id.storage_account_name_unique.hex
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}