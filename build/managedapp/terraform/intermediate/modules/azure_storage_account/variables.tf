variable "access_tier" {
  type        = string
  description = "Defines the access tier for BlobStorage, FileStorage and StorageV2 accounts."
  default     = "Hot"
}

variable "account_kind" {
  type        = string
  description = "Defines the Kind of account."
  default     = "StorageV2"
}

variable "account_replication_type" {
  type        = string
  description = "Defines the type of replication to use for this storage account."
}

variable "account_tier" {
  type        = string
  description = "Defines the Tier to use for this storage account."
}

variable "allow_nested_items_to_be_public" {
  default     = false
  description = <<-DESCRIPTION
  Whether to allow nested items within this account to opt into being public. This is
  the same setting in the portal called "Allow blob anonymous access".
  DESCRIPTION
  type        = bool
}

variable "containers" {
  type = map(object({
    container_access_type = optional(string, "private")
  }))
  description = "Container within an Azure Storage Account."
  default     = {}
}
variable "location" {
  type        = string
  description = "location"
  default     = "eastus2"
}

variable "monitor_diagnostic_destinations" {
  type = object({
    eventhubs = map(object({ # key = region
      authorization_rule_id = string
      eventhub_name         = string
      namespace_name        = string
    }))
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

variable "name" {
  type        = string
  description = "Name of the storage account"
}

variable "network_rules" {
  default = null
  type = object({
    bypass   = optional(set(string), ["AzureServices"])
    ip_rules = list(string)
  })
  description = <<-EOL
  Only one network rules block can be tied to a storage account. We currently don't have
  support added for `virtual_network_subnet_ids` or `private_link_access`. The `default_action`
  is automatically set to `Deny`. This action will be applied if there are no matching
  rules. `bypass` specifies whether traffic is bypassed for Logging/Metrics/AzureServices.
  Valid options are any combination of Logging, Metrics, AzureServices, or None. `ip_rules`
  should be a list of public IP or IP ranges in CIDR Format. Only IPv4 addresses are allowed.
  /31 CIDRs, /32 CIDRs, and Private IP address ranges are not allowed.
  EOL
}

variable "private_endpoints" {
  default     = {}
  description = <<-DESCRIPTION
  A collection of private endpoints to create & connect to the storage account. You
  must create a separate private endpoint per subresource. The VNet/subnet that is
  looked up must be in the same subscription as the storage account.
  DESCRIPTION
  type = map(object({
    subnet = object({
      name                 = string
      resource_group_name  = string
      virtual_network_name = string
    })
    subresource_name = string
  }))

  validation {
    condition = alltrue([for k, v in var.private_endpoints : contains([
      "blob",
      "dfs",
      "disks",
      "file",
      "queue",
      "table",
      "web",
    ], v.subresource_name)])
    error_message = "Must provide a valid subresource name"
  }

  validation {
    condition     = alltrue([for k, v in var.private_endpoints : strcontains(k, "pep-")])
    error_message = "Private endpoints should have the `pep-` prefix"
  }
}

variable "required_tags" {
  type = object({
    App         = string
    Environment = string
    GBU         = string
    ITSM        = optional(string, "STORAGE")
    JobWbs      = string
    Notes       = optional(string)
    Owner       = string
  })
  description = "Required Azure tags"

  validation {
    condition     = contains(["DEV", "DR", "PROD", "QA", "TEST"], var.required_tags.Environment)
    error_message = "Environment tag must be one of: DEV, TEST, QA, PROD, DR"
  }

  validation {
    condition     = contains(["COR", "FED", "INF", "MEA"], var.required_tags.GBU)
    error_message = "GBU tag must be one of: COR, FED, INF, MEA"
  }

  validation {
    condition = contains(
      ["BACKUP", "DATABASE", "MANAGEMENT", "NETWORK", "SERVER", "STORAGE"],
      var.required_tags.ITSM
    )
    error_message = "ITSM tag must be one of: BACKUP, DATABASE, MANAGEMENT, NETWORK, SERVER, STORAGE"
  }

  validation {
    condition     = can(regex("^\\d{6}-\\d{5}$", var.required_tags.JobWbs))
    error_message = "JobWbs tag must be digits in the format xxxxxx-xxxxx"
  }

  validation {
    condition     = can(regex("(?i)[A-Z0-9+_.-]+@[A-Z0-9.-]+", var.required_tags.Owner))
    error_message = "Owner tag must be a valid email address"
  }
}

variable "resource_group_name" {
  type        = string
  description = "resource group name"
}

variable "public_network_access_enabled" {
  type        = bool
  description = "Allow/Disallow public network access to the Storage Account"
  default     = false
}
