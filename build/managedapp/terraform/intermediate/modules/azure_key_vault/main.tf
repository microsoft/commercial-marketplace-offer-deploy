resource "azurerm_key_vault" "this" {
  enable_rbac_authorization       = var.enable_rbac_authorization
  enabled_for_disk_encryption     = var.enabled_for_disk_encryption
  enabled_for_template_deployment = var.enabled_for_template_deployment
  location                        = var.location
  name                            = var.key_vault_name
  public_network_access_enabled   = var.public_network_access_enabled
  purge_protection_enabled        = var.purge_protection_enabled
  resource_group_name             = var.resource_group_name
  sku_name                        = var.sku
  soft_delete_retention_days      = var.soft_delete_retention_days
  tags                            = var.required_tags
  tenant_id                       = data.azurerm_client_config.current.tenant_id
}

# Terraform operations fail without a role directly on the resource
resource "azurerm_role_assignment" "key_vault_admin" {
  for_each = (
    strcontains(var.location, "gov") ? local.gov_admin_principals :
    local.comm_admin_principals
  )

  principal_id         = each.value
  role_definition_name = "Key Vault Administrator"
  scope                = azurerm_key_vault.this.id
}

resource "azurerm_monitor_diagnostic_setting" "this" {
  eventhub_authorization_rule_id = var.monitor_diagnostic_destinations.eventhubs[var.location].authorization_rule_id
  eventhub_name                  = var.monitor_diagnostic_destinations.eventhubs[var.location].eventhub_name
  log_analytics_workspace_id     = var.monitor_diagnostic_destinations.log_analytics_workspace_id
  name                           = "diag-${var.key_vault_name}-activity-logs"
  target_resource_id             = azurerm_key_vault.this.id

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
    private_connection_resource_id = azurerm_key_vault.this.id
    subresource_names              = ["vault"]
  }

  lifecycle {
    # private_dns_zone_group gets configured automatically by Azure Policy
    ignore_changes = [private_dns_zone_group]
  }
}
