provider "azurerm" {
  features {}
}

variable "name" {
  description = "The name of the Linux web app"
  type        = string
}

variable "location" {
  description = "The location of the Linux web app"
  type        = string
}

variable "resourceGroupName" {
  description = "The resource group name for the Linux web app"
  type        = string
}

variable "custom_domains" {
  description = "The custom domains for the Linux web app"
  type        = list(string)
}

variable "enabled" {
  description = "Whether the Linux web app is enabled"
  type        = bool
}

variable "https_only" {
  description = "Whether to enforce HTTPS only"
  type        = bool
}

variable "service_plan_id" {
  description = "The service plan ID for the Linux web app"
  type        = string
}

variable "site_config" {
  description = "The site configuration for the Linux web app"
  type = object({
    always_on           = bool
    linux_fx_version    = string
    app_settings        = list(object({
      name  = string
      value = string
    }))
  })
}


variable "webdeploy_publish_basic_authentication_enabled" {
  description = "Whether to enable basic authentication for web deploy"
  type        = bool
}

variable "private_endpoints" {
  description = "The private endpoints for the Linux web app"
  type        = map(object({
    private_dns_zone_name = string
    subresource_names     = list(string)
    subnet                = object({
      name                 = string
      resource_group_name  = string
      virtual_network_name = string
    })
  }))
}

variable "monitor_diagnostic_destinations" {
  type = object({
    log_analytics_workspace_id = string
    resource_group_name        = string
    subscription_id            = string
  })
  description = <<-EOL
  Destinations used by azurerm_monitor_diagnostic_setting to store activity logs in a
  central location. The log analytics workspace doesn't have to be in the same region as
  the resource. The eventhub does have to be in the same region as the resource, so they
  are stored in a map where the key is the region.
  EOL
}

variable "app_service_plans" {
  description = "Map of app service plans"
  type        = map(object({
    maximum_elastic_worker_count = number
    os_type                      = string
    per_site_scaling_enabled     = bool
    sku_name                     = string
    worker_count                 = number
    zone_balancing_enabled       = bool
    tags                         = map(string)
  }))
}

variable "cognitive_accounts" {
  description = "Map of cognitive accounts"
  type        = map(object({
    resource_group_name                = string
    custom_subdomain_name              = string
    dynamic_throttling_enabled         = bool
    fqdns                              = list(string)
    kind                               = string
    local_auth_enabled                 = bool
    location                           = string
    outbound_network_access_restricted = bool
    public_network_access_enabled      = bool
    sku_name                           = string
    private_endpoints                  = map(object({
      private_dns_zone_name = string
      subresource_names     = list(string)
      subnet                = object({
        name                  = string
        resource_group_name   = string
        virtual_network_name  = string
      })
    }))
    tags = map(string)
  }))
}

variable "key_vaults" {
  description = "Map of key vaults"
  type        = map(object({
    resource_group_name                = string
    enable_rbac_authorization          = bool
    enabled_for_disk_encryption        = bool
    enabled_for_template_deployment    = bool
    location                           = string
    public_network_access_enabled      = bool
    purge_protection_enabled           = bool
    sku                                = string
    soft_delete_retention_days         = number
    private_endpoints                  = map(object({
      private_dns_zone_name = string
      subresource_names     = list(string)
      subnet                = object({
        name                  = string
        resource_group_name   = string
        virtual_network_name  = string
      })
    }))
    tags = map(string)
  }))
}

variable "linux_web_apps" {
  description = "Map of Linux web apps"
  type        = map(object({
    resource_group_name                = string
    custom_domains                     = list(string)
    enabled                            = bool
    https_only                         = bool
    location                           = string
    service_plan_key                    = string
    site_config                        = object({
      always_on            = bool
      linux_fx_version     = string
      app_settings         = list(object({
        name  = string
        value = string
      }))
    })
    webdeploy_publish_basic_authentication_enabled = bool
    private_endpoints                  = map(object({
      private_dns_zone_name = string
      subresource_names     = list(string)
      subnet                = object({
        name                  = string
        resource_group_name   = string
        virtual_network_name  = string
      })
    }))
    tags = map(string)
  }))
}

variable "postgresql_flexible_servers" {
  description = "Map of PostgreSQL flexible servers"
  type        = map(object({
    resource_group_name            = string
    admin_username                 = string
    auto_grow_enabled              = bool
    backup_retention_days          = number
    geo_redundant_backup_enabled   = bool
    location                       = string
    maintenance_window             = object({
      day_of_week  = number
      start_hour   = number
      start_minute = number
    })
    sku_name                       = string
    storage_mb                     = number
    storage_tier                   = string
    version                        = string
    vnet_integration               = object({
      private_dns_zone_id          = string
      delegated_subnet_name        = string
      vnet_name                    = string
      vnet_resource_group_name     = string
    })
    zone                           = number
    tags                           = map(string)
    high_availability              = object({
      mode                      = string
      standby_availability_zone = string
    })
  }))
}

variable "private_dns_zones" {
  description = "Map of private DNS zones"
  type        = map(object({
    resource_group_name     = string
    virtual_network_links   = map(object({
      virtual_network_id      = string
      registration_enabled    = bool
      resource_group_name     = string
    }))
    tags                    = map(string)
  }))
}

variable "search_services" {
  description = "Map of search services"
  type        = map(object({
    resource_group_name             = string
    authentication_failure_mode     = string
    hosting_mode                    = string
    local_authentication_enabled    = bool
    location                        = string
    public_network_access_enabled   = bool
    partition_count                 = number
    replica_count                   = number
    semantic_search_sku             = string
    sku                             = string
    private_endpoints               = map(object({
      private_dns_zone_name = string
      subresource_names     = list(string)
      subnet                = object({
        name                 = string
        resource_group_name  = string
        virtual_network_name = string
      })
    }))
    tags = map(string)
  }))
}

variable "storage_accounts" {
  description = "Map of storage accounts"
  type        = map(object({
    resource_group_name            = string
    access_tier                    = string
    account_kind                   = string
    account_replication_type       = string
    account_tier                   = string
    containers                     = map(object({
      container_access_type           = string
      encryption_scope_override_enabled = bool
      name                            = string
    }))
    location                       = string
    public_network_access_enabled  = bool
    network_rules                  = object({
      bypass         = list(string)
      default_action = string
      ip_rules       = list(string)
    })
    private_endpoints              = map(object({
      private_dns_zone_name = string
      subresource_name     = string
      subnet                = object({
        name                 = string
        resource_group_name  = string
        virtual_network_name = string
      })
    }))
    tags = map(string)
  }))
}

variable "default_tags" {
  description = "A map of default tags to apply to all resources"
  type        = map(string)
}


variable "required_tags" {
  type = object({
    App         = string
    Environment = string
    GBU         = string
    ITSM        = optional(string, "SERVER")
    JobWbs      = string
    Notes       = optional(string)
    Owner       = string
  })
  description = "Required Azure tags"
}

variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
}

variable "environment" {
  description = "The environment type, e.g., 'usgovernment' for government cloud or 'comm' for commercial cloud."
  type        = string
}

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
