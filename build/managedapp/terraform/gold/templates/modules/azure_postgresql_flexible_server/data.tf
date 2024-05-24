data "azurerm_monitor_diagnostic_categories" "this" {
  resource_id = azurerm_postgresql_flexible_server.this.id
}

data "azurerm_subnet" "this" {
  count = var.vnet_integration == null ? 0 : 1

  name                 = var.vnet_integration.delegated_subnet_name
  virtual_network_name = var.vnet_integration.vnet_name
  resource_group_name  = var.vnet_integration.vnet_resource_group_name
}
