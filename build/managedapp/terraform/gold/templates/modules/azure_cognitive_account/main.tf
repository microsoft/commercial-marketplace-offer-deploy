resource "azurerm_cognitive_account" "this" {
  custom_subdomain_name              = var.custom_subdomain_name
  dynamic_throttling_enabled         = var.dynamic_throttling_enabled
  fqdns                              = var.fqdns
  kind                               = var.kind
  local_auth_enabled                 = var.local_auth_enabled
  location                           = var.location
  name                               = var.name
  outbound_network_access_restricted = var.outbound_network_access_restricted
  public_network_access_enabled      = var.public_network_access_enabled
  resource_group_name                = var.resource_group_name
  sku_name                           = var.sku_name
  tags                               = var.required_tags
}

resource "azurerm_monitor_diagnostic_setting" "this" {
  log_analytics_workspace_id     = var.monitor_diagnostic_destinations.log_analytics_workspace_id
  name                           = "diag-${var.name}-activity-logs"
  target_resource_id             = azurerm_cognitive_account.this.id

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

resource "azurerm_private_endpoint" "this" {
  for_each = var.private_endpoints

  custom_network_interface_name = "nic-${each.key}"
  location                      = var.location
  name                          = each.key
  resource_group_name           = var.resource_group_name
  subnet_id                     = data.azurerm_subnet.private_endpoint[each.key].id
  tags                          = var.required_tags

  private_service_connection {
    is_manual_connection           = false
    name                           = each.key
    private_connection_resource_id = azurerm_cognitive_account.this.id
    subresource_names              = ["account"]
  }

  lifecycle {
    # private_dns_zone_group gets configured automatically by Azure Policy
    ignore_changes = [private_dns_zone_group]
  }
}
