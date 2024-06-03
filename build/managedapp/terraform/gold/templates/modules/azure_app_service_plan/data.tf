data "azurerm_monitor_diagnostic_categories" "this" {
  resource_id = azurerm_service_plan.this.id
}
