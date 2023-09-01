variable "location" {}
variable "resource_group_name" {}

resource "azurerm_storage_account" "example" {
  name                     = "bobjacmodmterra"
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
