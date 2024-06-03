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