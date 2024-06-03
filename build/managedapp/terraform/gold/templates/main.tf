provider "azurerm" {
  features {}
}

module "azure_app_service_plan" {
  source                         = "./modules/azure_app_service_plan"
  for_each                       = var.app_service_plans
  name                           = each.key
  location                       = var.location
  resource_group_name            = var.resourceGroupName
  maximum_elastic_worker_count   = each.value.maximum_elastic_worker_count
  os_type                        = each.value.os_type
  per_site_scaling_enabled       = each.value.per_site_scaling_enabled
  sku_name                       = each.value.sku_name
  worker_count                   = each.value.worker_count
  zone_balancing_enabled         = each.value.zone_balancing_enabled
  monitor_diagnostic_destinations = var.monitor_diagnostic_destinations
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_cognitive_account" {
  source                         = "./modules/azure_cognitive_account"
  for_each                       = var.cognitive_accounts
  name                           = each.key
  location                       = each.value.location
  resource_group_name            = var.resourceGroupName
  custom_subdomain_name          = each.value.custom_subdomain_name
  dynamic_throttling_enabled     = each.value.dynamic_throttling_enabled
  fqdns                          = each.value.fqdns
  kind                           = each.value.kind
  local_auth_enabled             = each.value.local_auth_enabled
  outbound_network_access_restricted = each.value.outbound_network_access_restricted
  public_network_access_enabled  = each.value.public_network_access_enabled
  sku_name                       = each.value.sku_name
  private_endpoints              = each.value.private_endpoints
  monitor_diagnostic_destinations = var.monitor_diagnostic_destinations
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_key_vault" {
  source                         = "./modules/azure_key_vault"
  for_each                       = var.key_vaults
  key_vault_name                 = each.key
  location                       = each.value.location
  resource_group_name            = var.resourceGroupName
  enable_rbac_authorization      = each.value.enable_rbac_authorization
  enabled_for_disk_encryption    = each.value.enabled_for_disk_encryption
  enabled_for_template_deployment = each.value.enabled_for_template_deployment
  public_network_access_enabled  = each.value.public_network_access_enabled
  purge_protection_enabled       = each.value.purge_protection_enabled
  sku                            = each.value.sku
  soft_delete_retention_days     = each.value.soft_delete_retention_days
  private_endpoints              = each.value.private_endpoints
  monitor_diagnostic_destinations = var.monitor_diagnostic_destinations
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_linux_web_app" {
  source                         = "./modules/azure_linux_web_app"
  for_each                       = var.linux_web_apps
  name                           = each.key
  location                       = each.value.location
  resource_group_name            = var.resourceGroupName
  custom_domains                 = each.value.custom_domains
  enabled                        = each.value.enabled
  https_only                     = each.value.https_only
  service_plan_id                = module.azure_app_service_plan[each.value.service_plan_key].id
  site_config                    = each.value.site_config
  webdeploy_publish_basic_authentication_enabled = each.value.webdeploy_publish_basic_authentication_enabled
  private_endpoints              = each.value.private_endpoints
  monitor_diagnostic_destinations = var.monitor_diagnostic_destinations
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_postgresql_flexible_server" {
  source                         = "./modules/azure_postgresql_flexible_server"
  for_each                       = var.postgresql_flexible_servers
  name                           = each.key
  is_gov                         = false
  location                       = each.value.location
  resource_group_name            = var.resourceGroupName
  postgresql_version             = each.value.version
  admin_username                 = each.value.admin_username
  auto_grow_enabled              = each.value.auto_grow_enabled
  backup_retention_days          = each.value.backup_retention_days
  geo_redundant_backup_enabled   = each.value.geo_redundant_backup_enabled
  maintenance_window             = each.value.maintenance_window
  sku_name                       = each.value.sku_name
  storage_mb                     = each.value.storage_mb
  storage_tier                   = each.value.storage_tier
  vnet_integration               = each.value.vnet_integration
  zone                           = each.value.zone
  monitor_diagnostic_destinations = var.monitor_diagnostic_destinations
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_private_dns_zone" {
  source                         = "./modules/azure_private_dns_zone"
  for_each                       = var.private_dns_zones
  name                           = each.key
  resource_group_name            = var.resourceGroupName
  virtual_network_links          = each.value.virtual_network_links
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_search_service" {
  source                         = "./modules/azure_search_service"
  for_each                       = var.search_services
  name                           = each.key
  location                       = each.value.location
  resource_group_name            = var.resourceGroupName
  authentication_failure_mode    = each.value.authentication_failure_mode
  hosting_mode                   = each.value.hosting_mode
  local_authentication_enabled   = each.value.local_authentication_enabled
  public_network_access_enabled  = each.value.public_network_access_enabled
  partition_count                = each.value.partition_count
  replica_count                  = each.value.replica_count
  semantic_search_sku            = each.value.semantic_search_sku
  sku                            = each.value.sku
  private_endpoints              = each.value.private_endpoints
  monitor_diagnostic_destinations = var.monitor_diagnostic_destinations
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_storage_account" {
  source                         = "./modules/azure_storage_account"
  for_each                       = var.storage_accounts
  name                           = each.key
  location                       = each.value.location
  resource_group_name            = var.resourceGroupName
  access_tier                    = each.value.access_tier
  account_kind                   = each.value.account_kind
  account_replication_type       = each.value.account_replication_type
  account_tier                   = each.value.account_tier
  containers                     = each.value.containers
  public_network_access_enabled  = each.value.public_network_access_enabled
  network_rules                  = each.value.network_rules
  private_endpoints              = each.value.private_endpoints
  monitor_diagnostic_destinations = var.monitor_diagnostic_destinations
  required_tags                  = merge(var.default_tags, each.value.tags)
}

module "azure_resource_provider_registration" {
  source                         = "./modules/azure_resource_provider_registration"
  for_each                       = local.resource_providers_to_register
  name                           = each.value
  subscription_id                = var.subscription_id
  platform                       = "linux"
}
