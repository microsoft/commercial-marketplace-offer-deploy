resource "azurerm_search_service" "this" {
  authentication_failure_mode   = var.authentication_failure_mode
  hosting_mode                  = var.hosting_mode
  local_authentication_enabled  = var.local_authentication_enabled
  location                      = var.location
  name                          = var.name
  partition_count               = var.partition_count
  public_network_access_enabled = var.public_network_access_enabled
  replica_count                 = var.replica_count
  resource_group_name           = var.resource_group_name
  tags                          = var.required_tags
  semantic_search_sku           = var.semantic_search_sku
  sku                           = var.sku

  lifecycle {
    precondition {
      condition     = !(var.authentication_failure_mode != null && var.local_authentication_enabled == false)
      error_message = "`authentication_failure_mode` can only be set when `local_authentication_enabled` is true"
    }
    precondition {
      condition     = !(var.hosting_mode != "default" && var.sku != "standard3")
      error_message = "`hosting_mode` can only be set when `sku` is set to `standard3`"
    }
    precondition {
      condition     = !(var.partition_count != 1 && (var.sku == "free" || var.sku == "basic"))
      error_message = "`partition_count` can't be set when `sku` is set to `free` or `basic`"
    }
    precondition {
      condition     = !(var.semantic_search_sku != null && var.sku == "free")
      error_message = "`semantic_search_sku` can't be set when `sku` is set to `free`"
    }
    precondition {
      condition     = !(var.semantic_search_sku != null && strcontains(var.location, "gov"))
      error_message = "`semantic_search_sku` isn't supported in Azure Government"
    }
  }
}

resource "azurerm_monitor_diagnostic_setting" "this" {
  log_analytics_workspace_id     = var.monitor_diagnostic_destinations.log_analytics_workspace_id
  name                           = "diag-${var.name}-activity-logs"
  target_resource_id             = azurerm_search_service.this.id

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
    private_connection_resource_id = azurerm_search_service.this.id
    subresource_names              = ["searchService"]
  }

  lifecycle {
    # private_dns_zone_group gets configured automatically by Azure Policy
    ignore_changes = [private_dns_zone_group]
  }
}
