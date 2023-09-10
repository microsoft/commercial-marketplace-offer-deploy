

data "azurerm_storage_account" "example" {
  name                     = "bobjacterra3"
  resource_group_name      = data.azurerm_resource_group.example.name 
}
