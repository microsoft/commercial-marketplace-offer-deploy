resource "random_password" "this" {
  length = 16
}

resource "azurerm_postgresql_flexible_server" "this" {
  administrator_login          = var.admin_username
  administrator_password       = random_password.this.result
  auto_grow_enabled            = var.auto_grow_enabled
  backup_retention_days        = var.backup_retention_days
  delegated_subnet_id          = try(data.azurerm_subnet.this[0].id, null)
  geo_redundant_backup_enabled = var.geo_redundant_backup_enabled
  location                     = var.location
  name                         = var.name
  private_dns_zone_id          = var.vnet_integration == null ? null : local.private_dns_zone_id
  resource_group_name          = var.resource_group_name
  sku_name                     = var.sku_name
  storage_mb                   = var.storage_mb
  storage_tier                 = var.storage_tier
  tags                         = var.required_tags
  version                      = var.postgresql_version
  zone                         = var.zone

  dynamic "high_availability" {
    for_each = var.high_availability == null ? [] : [1]
    content {
      mode                      = var.high_availability.mode
      standby_availability_zone = var.high_availability.standby_availability_zone
    }
  }

  dynamic "maintenance_window" {
    for_each = var.maintenance_window == null ? [] : [1]
    content {
      day_of_week  = var.maintenance_window.day_of_week
      start_hour   = var.maintenance_window.start_hour
      start_minute = var.maintenance_window.start_minute
    }
  }

  lifecycle {
    ignore_changes = [
      high_availability["standby_availability_zone"],
      zone,
    ]

    precondition {
      condition     = !(startswith(var.sku_name, "B_") && var.high_availability != null)
      error_message = "Burstable SKU's don't support high availability"
    }
  }
}

resource "azurerm_monitor_diagnostic_setting" "this" {
  eventhub_authorization_rule_id = var.monitor_diagnostic_destinations.eventhubs[var.location].authorization_rule_id
  eventhub_name                  = var.monitor_diagnostic_destinations.eventhubs[var.location].eventhub_name
  log_analytics_workspace_id     = var.monitor_diagnostic_destinations.log_analytics_workspace_id
  name                           = "diag-${var.name}-activity-logs"
  target_resource_id             = azurerm_postgresql_flexible_server.this.id

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
