locals {
  comm_monitor_diagnostic_destinations = {
    "eventhubs" = {
      "eastus" = {
        "authorization_rule_id" = ""
        "eventhub_name"         = ""
        "namespace_name"        = ""
      }
      "westus" = {
        "authorization_rule_id" = ""
        "eventhub_name"         = ""
        "namespace_name"        = ""
      }
    }
    "log_analytics_workspace_id" = ""
    "resource_group_name"        = ""
    "subscription_id"            = ""
  }
  gov_monitor_diagnostic_destinations = {
    "eventhubs" = {
      "usgovvirginia" = {
        "authorization_rule_id" = ""
        "eventhub_name"         = ""
        "namespace_name"        = ""
      }
    }
    "log_analytics_workspace_id" = ""
    "resource_group_name"        = ""
    "subscription_id"            = ""
  }
  is_gov = var.environment == "usgovernment"
  monitor_diagnostic_destinations = (
    local.is_gov ? local.gov_monitor_diagnostic_destinations : local.comm_monitor_diagnostic_destinations
  )
  resource_providers_to_register = toset([for item in [
    length(var.app_service_plans) > 0 ? "Microsoft.Web" : null,
    length(var.cognitive_accounts) > 0 ? "Microsoft.CognitiveServices" : null,
    length(var.cognitive_accounts) > 0 ? "Microsoft.Network" : null,
    length(var.key_vaults) > 0 ? "Microsoft.KeyVault" : null,
    length(var.linux_web_apps) > 0 ? "Microsoft.Network" : null,
    length(var.linux_web_apps) > 0 ? "Microsoft.Web" : null,
    length(var.postgresql_flexible_servers) > 0 ? "Microsoft.DBforPostgreSQL" : null,
    length(var.private_dns_zones) > 0 ? "Microsoft.Network" : null,
    length(var.search_services) > 0 ? "Microsoft.Network" : null,
    length(var.search_services) > 0 ? "Microsoft.Search" : null,
    length(var.storage_accounts) > 0 ? "Microsoft.Network" : null,
    length(var.storage_accounts) > 0 ? "Microsoft.Storage" : null,
  ] : item if item != null])
}
