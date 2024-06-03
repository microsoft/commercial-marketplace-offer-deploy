resource "azurerm_service_plan" "this" {
  location                     = var.location
  maximum_elastic_worker_count = var.maximum_elastic_worker_count
  name                         = var.name
  os_type                      = var.os_type
  per_site_scaling_enabled     = var.per_site_scaling_enabled
  resource_group_name          = var.resource_group_name
  sku_name                     = var.sku_name
  tags                         = var.required_tags
  worker_count                 = var.worker_count
  zone_balancing_enabled       = var.zone_balancing_enabled
}

resource "azurerm_monitor_diagnostic_setting" "this" {
  log_analytics_workspace_id     = var.monitor_diagnostic_destinations.log_analytics_workspace_id
  name                           = "diag-${var.name}-activity-logs"
  target_resource_id             = azurerm_service_plan.this.id

  dynamic "enabled_log" {
    for_each = toset(data.azurerm_monitor_diagnostic_categories.this.log_category_types)
    content {
      category = enabled_log.key
    }
  }

  dynamic "metric" {
    for_each = toset(data.azurerm_monitor_diagnostic_categories.this.metrics)
    content {
      category = metric.key
      enabled  = true
    }
  }
}
