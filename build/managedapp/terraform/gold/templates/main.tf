module "azure_app_service_plan" {
  for_each = var.app_service_plans
  source   = "./modules/azure_app_service_plan"

  location                        = var.location
  maximum_elastic_worker_count    = each.value.maximum_elastic_worker_count
  monitor_diagnostic_destinations = local.monitor_diagnostic_destinations
  name                            = each.key
  os_type                         = each.value.os_type
  per_site_scaling_enabled        = each.value.per_site_scaling_enabled
  required_tags                   = merge(var.default_tags, each.value.tags)
  resource_group_name             = var.resourceGroupName
  sku_name               = each.value.sku_name
  worker_count           = each.value.worker_count
  zone_balancing_enabled = each.value.zone_balancing_enabled

  depends_on = [module.azure_resource_provider_registration["Microsoft.Web"]]
}

module "azure_cognitive_account" {
  for_each = var.cognitive_accounts
  source   = "./modules/azure_cognitive_account"

  custom_subdomain_name              = each.value.custom_subdomain_name
  dynamic_throttling_enabled         = each.value.dynamic_throttling_enabled
  fqdns                              = each.value.fqdns
  kind                               = each.value.kind
  local_auth_enabled                 = each.value.local_auth_enabled
  location                           = var.location
  monitor_diagnostic_destinations    = local.monitor_diagnostic_destinations
  name                               = each.key
  outbound_network_access_restricted = each.value.outbound_network_access_restricted
  private_endpoints                  = each.value.private_endpoints
  public_network_access_enabled      = each.value.public_network_access_enabled
  required_tags                      = merge(var.default_tags, each.value.tags)
  resource_group_name                = var.resourceGroupName
  sku_name = each.value.sku_name

  depends_on = [
    module.azure_resource_provider_registration["Microsoft.CognitiveServices"],
    module.azure_resource_provider_registration["Microsoft.Network"],
  ]
}

module "azure_key_vault" {
  for_each = var.key_vaults
  source   = "./modules/azure_key_vault"

  enable_rbac_authorization       = each.value.enable_rbac_authorization
  enabled_for_disk_encryption     = each.value.enabled_for_disk_encryption
  enabled_for_template_deployment = each.value.enabled_for_template_deployment
  key_vault_name                  = each.key
  location                        = var.location
  monitor_diagnostic_destinations = local.monitor_diagnostic_destinations
  private_endpoints               = each.value.private_endpoints
  public_network_access_enabled   = each.value.public_network_access_enabled
  purge_protection_enabled        = each.value.purge_protection_enabled
  required_tags                   = merge(var.default_tags, each.value.tags)
  resource_group_name            = var.resourceGroupName
  sku                        = each.value.sku
  soft_delete_retention_days = each.value.soft_delete_retention_days

  depends_on = [module.azure_resource_provider_registration["Microsoft.KeyVault"]]
}

module "azure_linux_web_app" {
  for_each = var.linux_web_apps
  source   = "./modules/azure_linux_web_app"

  custom_domains                  = each.value.custom_domains
  enabled                         = each.value.enabled
  https_only                      = each.value.https_only
  location                        = var.location
  monitor_diagnostic_destinations = local.monitor_diagnostic_destinations
  name                            = each.key
  private_endpoints               = each.value.private_endpoints
  public_network_access_enabled   = each.value.public_network_access_enabled
  required_tags                   = merge(var.default_tags, each.value.tags)
  resource_group_name             = var.resourceGroupName
  service_plan_id                                = module.azure_app_service_plan[each.value.service_plan_key].id
  site_config                                    = each.value.site_config
  webdeploy_publish_basic_authentication_enabled = each.value.webdeploy_publish_basic_authentication_enabled

  depends_on = [
    module.azure_resource_provider_registration["Microsoft.Network"],
    module.azure_resource_provider_registration["Microsoft.Web"],
  ]
}

module "azure_postgresql_flexible_server" {
  for_each = var.postgresql_flexible_servers
  source   = "./modules/azure_postgresql_flexible_server"

  admin_username                  = each.value.admin_username
  auto_grow_enabled               = each.value.auto_grow_enabled
  backup_retention_days           = each.value.backup_retention_days
  geo_redundant_backup_enabled    = each.value.geo_redundant_backup_enabled
  high_availability               = each.value.high_availability
  is_gov                          = local.is_gov
  location                        = var.location
  maintenance_window              = each.value.maintenance_window
  monitor_diagnostic_destinations = local.monitor_diagnostic_destinations
  name                            = each.key
  postgresql_version              = each.value.version
  required_tags                   = merge(var.default_tags, each.value.tags)
  resource_group_name             = var.resourceGroupName
  sku_name         = each.value.sku_name
  storage_mb       = each.value.storage_mb
  storage_tier     = each.value.storage_tier
  vnet_integration = each.value.vnet_integration
  zone             = each.value.zone

  depends_on = [module.azure_resource_provider_registration["Microsoft.DBforPostgreSQL"]]
}

module "azure_private_dns_zone" {
  for_each = var.private_dns_zones
  source   = "./modules/azure_private_dns_zone"

  name          = each.key
  required_tags = merge(var.default_tags, each.value.tags)
  resource_group_name = var.resourceGroupName
  virtual_network_links = each.value.virtual_network_links

  depends_on = [module.azure_resource_provider_registration["Microsoft.Network"]]
}

module "azure_resource_provider_registration" {
  for_each = local.resource_providers_to_register
  source   = "./modules/azure_resource_provider_registration"

  name            = each.value
  platform        = var.platform
  subscription_id = var.subscription_id
}

module "azure_search_service" {
  for_each = var.search_services
  source   = "./modules/azure_search_service"

  authentication_failure_mode     = each.value.authentication_failure_mode
  hosting_mode                    = each.value.hosting_mode
  local_authentication_enabled    = each.value.local_authentication_enabled
  location                        = var.location
  monitor_diagnostic_destinations = local.monitor_diagnostic_destinations
  name                            = each.key
  partition_count                 = each.value.partition_count
  private_endpoints               = each.value.private_endpoints
  public_network_access_enabled   = each.value.public_network_access_enabled
  replica_count                   = each.value.replica_count
  required_tags                   = merge(var.default_tags, each.value.tags)
  resource_group_name             = var.resourceGroupName
  semantic_search_sku = each.value.semantic_search_sku
  sku                 = each.value.sku

  depends_on = [
    module.azure_resource_provider_registration["Microsoft.Network"],
    module.azure_resource_provider_registration["Microsoft.Search"],
  ]
}

module "azure_storage_account" {
  for_each = var.storage_accounts
  source   = "./modules/azure_storage_account"

  access_tier                     = each.value.access_tier
  account_kind                    = each.value.account_kind
  account_replication_type        = each.value.account_replication_type
  account_tier                    = each.value.account_tier
  containers                      = each.value.containers
  location                        = var.location
  monitor_diagnostic_destinations = local.monitor_diagnostic_destinations
  name                            = each.key
  network_rules                   = each.value.network_rules
  private_endpoints               = each.value.private_endpoints
  public_network_access_enabled   = each.value.public_network_access_enabled
  required_tags                   = merge(var.default_tags, each.value.tags)
  resource_group_name             = var.resourceGroupName

  depends_on = [
    module.azure_resource_provider_registration["Microsoft.Network"],
    module.azure_resource_provider_registration["Microsoft.Storage"],
  ]
}
