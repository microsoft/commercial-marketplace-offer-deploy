data "azurerm_monitor_diagnostic_categories" "this" {
  resource_id = azurerm_storage_account.this.id
}

data "azurerm_subnet" "private_endpoint" {
  for_each = var.private_endpoints

  name                 = each.value.subnet.name
  resource_group_name  = each.value.subnet.resource_group_name
  virtual_network_name = each.value.subnet.virtual_network_name
}
