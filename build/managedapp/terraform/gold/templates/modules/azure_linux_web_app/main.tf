resource "azurerm_linux_web_app" "this" {
  enabled                                        = var.enabled
  name                                           = var.name
  resource_group_name                            = var.resource_group_name
  location                                       = var.location
  service_plan_id                                = var.service_plan_id
  https_only                                     = var.https_only
  public_network_access_enabled                  = var.public_network_access_enabled
  tags                                           = var.required_tags
  webdeploy_publish_basic_authentication_enabled = var.webdeploy_publish_basic_authentication_enabled

  site_config {
    always_on                         = var.site_config.always_on
    ftps_state                        = var.site_config.ftps_state
    http2_enabled                     = var.site_config.http2_enabled
    ip_restriction_default_action     = var.site_config.ip_restriction_default_action
    load_balancing_mode               = var.site_config.load_balancing_mode
    minimum_tls_version               = var.site_config.minimum_tls_version
    scm_ip_restriction_default_action = var.site_config.scm_ip_restriction_default_action
    use_32_bit_worker                 = var.site_config.use_32_bit_worker
    vnet_route_all_enabled            = var.site_config.vnet_route_all_enabled

    dynamic "application_stack" {
      for_each = var.site_config.application_stack == null ? [] : [1]
      content {
        dotnet_version      = var.site_config.application_stack.dotnet_version
        go_version          = var.site_config.application_stack.go_version
        java_server         = var.site_config.application_stack.java_server
        java_server_version = var.site_config.application_stack.java_server_version
        java_version        = var.site_config.application_stack.java_version
        node_version        = var.site_config.application_stack.node_version
        php_version         = var.site_config.application_stack.php_version
        python_version      = var.site_config.application_stack.python_version
        ruby_version        = var.site_config.application_stack.ruby_version
      }
    }
  }

  lifecycle {
    ignore_changes = [app_settings]
  }
}

resource "azurerm_monitor_diagnostic_setting" "this" {
  log_analytics_workspace_id     = var.monitor_diagnostic_destinations.log_analytics_workspace_id
  name                           = "diag-${var.name}-activity-logs"
  target_resource_id             = azurerm_linux_web_app.this.id

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
    private_connection_resource_id = azurerm_linux_web_app.this.id
    subresource_names              = ["sites"]
  }

  lifecycle {
    # private_dns_zone_group gets configured automatically by Azure Policy
    ignore_changes = [private_dns_zone_group]
  }
}

resource "azurerm_app_service_custom_hostname_binding" "this" {
  for_each = var.custom_domains

  hostname            = each.value
  app_service_name    = azurerm_linux_web_app.this.name
  resource_group_name = var.resource_group_name
}

resource "azurerm_app_service_managed_certificate" "this" {
  for_each = var.custom_domains

  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.this[each.key].id
  tags                       = var.required_tags
}

resource "azurerm_app_service_certificate_binding" "this" {
  for_each            = var.custom_domains
  hostname_binding_id = azurerm_app_service_custom_hostname_binding.this[each.key].id
  certificate_id      = azurerm_app_service_managed_certificate.this[each.key].id
  ssl_state           = "SniEnabled"
}
