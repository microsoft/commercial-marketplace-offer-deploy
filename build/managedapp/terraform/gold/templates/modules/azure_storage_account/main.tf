resource "azurerm_storage_account" "this" {
  access_tier                     = var.access_tier
  account_kind                    = var.account_kind
  account_replication_type        = var.account_replication_type
  account_tier                    = var.account_tier
  allow_nested_items_to_be_public = var.allow_nested_items_to_be_public
  enable_https_traffic_only       = true
  location                        = var.location
  min_tls_version                 = "TLS1_2"
  name                            = var.name
  public_network_access_enabled   = var.public_network_access_enabled
  resource_group_name             = var.resource_group_name
  tags                            = var.required_tags

  dynamic "network_rules" {
    for_each = var.network_rules == null ? {} : { network_rules = var.network_rules }
    content {
      bypass         = var.network_rules.bypass
      default_action = "Deny"
      ip_rules       = var.network_rules.ip_rules
    }
  }
}

resource "azurerm_storage_container" "this" {
  for_each = var.containers

  container_access_type = each.value.container_access_type
  name                  = each.key
  storage_account_name  = azurerm_storage_account.this.name
}

resource "azurerm_monitor_diagnostic_setting" "this" {
  eventhub_authorization_rule_id = var.monitor_diagnostic_destinations.eventhubs[var.location].authorization_rule_id
  eventhub_name                  = var.monitor_diagnostic_destinations.eventhubs[var.location].eventhub_name
  log_analytics_workspace_id     = var.monitor_diagnostic_destinations.log_analytics_workspace_id
  name                           = "diag-${var.name}-activity-logs"
  target_resource_id             = azurerm_storage_account.this.id
  ##log_analytics_destination_type = var.diag_destination_type "unsure if we use this"

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
      enabled  = true # Create other Modules if you don't want them all on, by default
    }
  }
}

resource "azurerm_private_endpoint" "this" {
  for_each = var.private_endpoints

  custom_network_interface_name = "nic-${each.key}"
  name                          = each.key
  location                      = var.location
  resource_group_name           = var.resource_group_name
  subnet_id                     = data.azurerm_subnet.private_endpoint[each.key].id
  tags                          = var.required_tags

  private_service_connection {
    is_manual_connection           = false
    name                           = each.key
    private_connection_resource_id = azurerm_storage_account.this.id
    subresource_names              = [each.value.subresource_name]
  }

  lifecycle {
    # private_dns_zone_group gets configured automatically by Azure Policy
    ignore_changes = [private_dns_zone_group]
  }
}
