variable "app_service_plans" {
  default     = {}
  description = "A collection of App Service Service Plans to deploy"
  type = map(object({
    location                     = string
    maximum_elastic_worker_count = optional(number)
    os_type                      = string
    per_site_scaling_enabled     = optional(bool, false)
    resource_group_name          = string
    sku_name                     = string
    tags                         = optional(map(string))
    worker_count                 = optional(number)
    zone_balancing_enabled       = optional(bool)
  }))
}

variable "cognitive_accounts" {
  default     = {}
  description = "A map of Azure cognitive accounts to create"
  type = map(object({
    custom_subdomain_name              = string
    dynamic_throttling_enabled         = optional(bool)
    fqdns                              = optional(list(string))
    kind                               = string
    local_auth_enabled                 = optional(bool, true)
    location                           = string
    outbound_network_access_restricted = optional(bool, false)
    private_endpoints = optional(map(object({
      subnet = object({
        name                 = string
        resource_group_name  = string
        virtual_network_name = string
      })
    })), {})
    public_network_access_enabled = optional(bool, false)
    resource_group_name           = string
    sku_name                      = string
    tags                          = optional(map(string))
  }))
}

variable "default_tags" {
  type = object({
    App         = string
    Environment = string
    GBU         = string
    JobWbs      = string
    Notes       = optional(string)
    Owner       = string
  })
  description = "Required default tags which can be overriden"
}

variable "environment" {
  type        = string
  default     = null
  description = "The cloud env which should be used. Used to set cloud to usgovernment"
}

variable "key_vaults" {
  type = map(object({
    enable_rbac_authorization       = optional(bool, true)
    enabled_for_disk_encryption     = optional(bool)
    enabled_for_template_deployment = optional(bool)
    location                        = string
    private_endpoints = optional(map(object({
      subnet = object({
        name                 = string
        resource_group_name  = string
        virtual_network_name = string
      })
    })), {})
    public_network_access_enabled = optional(bool, false)
    purge_protection_enabled      = optional(bool)
    resource_group_name           = string
    sku                           = optional(string, "standard")
    soft_delete_retention_days    = optional(number, 90)
    tags                          = optional(map(string))
  }))
  description = "Map of objects for azure key vaults"
  default     = {}
}

variable "linux_web_apps" {
  default     = {}
  description = "A collection of Linux web apps to deploy"
  type = map(object({
    custom_domains = optional(set(string), [])
    enabled        = optional(bool, true)
    https_only     = optional(bool, false)
    location       = string
    private_endpoints = optional(map(object({
      subnet = object({
        name                 = string
        resource_group_name  = string
        virtual_network_name = string
      })
    })), {})
    public_network_access_enabled = optional(bool, false)
    resource_group_name           = string
    service_plan_key              = string
    site_config = object({
      always_on = optional(bool, true)
      application_stack = optional(object({
        dotnet_version      = optional(string)
        go_version          = optional(string)
        java_server         = optional(string)
        java_server_version = optional(string)
        java_version        = optional(number)
        node_version        = optional(string)
        php_version         = optional(number)
        python_version      = optional(number)
        ruby_version        = optional(number)
      }))
      ftps_state                        = optional(string, "Disabled")
      http2_enabled                     = optional(bool)
      ip_restriction_default_action     = optional(string, "Allow")
      load_balancing_mode               = optional(string, "LeastRequests")
      minimum_tls_version               = optional(number, 1.2)
      scm_ip_restriction_default_action = optional(string, "Allow")
      use_32_bit_worker                 = optional(bool, true)
      vnet_route_all_enabled            = optional(bool, false)
    })
    tags                                           = optional(map(string))
    webdeploy_publish_basic_authentication_enabled = optional(bool, true)
  }))
}

variable "platform" {
  type        = string
  description = <<-EOL
  The current OS, determined by Terragrunt. Some of the returned values can be:
  darwin, freebsd, linux, windows
  See: https://terragrunt.gruntwork.io/docs/reference/built-in-functions/#get_platform
  EOL
}

variable "postgresql_flexible_servers" {
  default     = {}
  description = "A collection of PostgreSql flexible servers to deploy"
  type = map(object({
    admin_username               = string
    auto_grow_enabled            = optional(bool, false)
    backup_retention_days        = optional(number, 7)
    geo_redundant_backup_enabled = optional(bool, false)
    high_availability = optional(object({
      mode                      = string
      standby_availability_zone = optional(number)
    }))
    location = string
    maintenance_window = optional(object({
      day_of_week  = optional(number, 0)
      start_hour   = optional(number, 0)
      start_minute = optional(number, 0)
    }))
    resource_group_name = string
    sku_name            = string
    storage_mb          = optional(number, 32768)
    storage_tier        = optional(string)
    tags                = optional(map(string))
    version             = number
    vnet_integration = optional(object({
      delegated_subnet_name    = string
      vnet_name                = string
      vnet_resource_group_name = string
    }))
    zone = optional(number)
  }))
}

variable "private_dns_zones" {
  default     = {}
  description = "value"
  type = map(object({
    resource_group_name = string
    tags                = optional(map(string))
    virtual_network_links = optional(map(object({ # key = VNet name
      auto_registration_enabled = optional(bool, false)
      resource_group_name       = string
    })))
  }))
}

variable "resource_groups" {
  type        = map(string)
  description = "Resource groups to create. The key is the name and the value is the region"
  default     = {}
}

variable "search_services" {
  default     = {}
  description = "A collection of Azure AI Search Services to deploy"
  type = map(object({
    authentication_failure_mode  = optional(string)
    hosting_mode                 = optional(string, "default")
    local_authentication_enabled = optional(bool, true)
    location                     = string
    partition_count              = optional(number, 1)
    private_endpoints = optional(map(object({
      subnet = object({
        name                 = string
        resource_group_name  = string
        virtual_network_name = string
      })
    })), {})
    public_network_access_enabled = optional(bool, false)
    replica_count                 = optional(number)
    resource_group_name           = string
    semantic_search_sku           = optional(string)
    sku                           = string
    tags                          = optional(map(string))
  }))
}

variable "storage_accounts" {
  type = map(object({
    access_tier              = optional(string, "Hot")
    account_kind             = optional(string, "StorageV2")
    account_replication_type = string
    account_tier             = string
    containers = optional(map(object({
      container_access_type = optional(string, "private")
    })), {})
    location = string
    network_rules = optional(object({
      bypass   = optional(set(string), ["AzureServices"])
      ip_rules = list(string)
    }))
    private_endpoints = optional(map(object({
      subnet = object({
        name                 = string
        resource_group_name  = string
        virtual_network_name = string
      })
      subresource_name = string
    })), {})
    public_network_access_enabled = optional(bool, false)
    resource_group_name           = string
    tags                          = optional(map(string))
  }))
  description = "Map of objects for Storage Accounts"
  default     = {}
}

variable "subscription_id" {
  type        = string
  description = "ID for the Azure subscription"
  default     = null

  validation {
    condition     = var.subscription_id == lower(var.subscription_id)
    error_message = "The subscription ID should be in all lower case."
  }
}
