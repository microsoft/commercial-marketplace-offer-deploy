variable "enable_rbac_authorization" {
  default     = true
  description = <<-DESCRIPTION
  Specify whether Azure Key Vault uses Role Based Access Control (RBAC) for authorization
  of data actions.
  DESCRIPTION
  type        = bool
}

variable "enabled_for_disk_encryption" {
  default     = null
  description = <<-DESCRIPTION
  Specify whether Azure Disk Encryption is permitted to retrieve secrets from the vault
  and unwrap keys.
  DESCRIPTION
  type        = bool
}

variable "enabled_for_template_deployment" {
  default     = null
  description = <<-DESCRIPTION
  Specify whether Azure Resource Manager is permitted to retrieve secrets from the key
  vault.
  DESCRIPTION
  type        = bool
}

variable "key_vault_name" {
  type        = string
  description = "Name of key vault to be created"
}

variable "location" {
  type        = string
  description = "Azure location"
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

variable "private_endpoints" {
  default     = {}
  description = <<-DESCRIPTION
  A collection of private endpoints to create & connect to the key vault.
  `subresource_names` defaults to `['vault']`, as that is the only available
  subresource. The VNet/subnet that is looked up must be in the same subscription as the
  key vault.
  DESCRIPTION
  type = map(object({
    subnet = object({
      name                 = string
      resource_group_name  = string
      virtual_network_name = string
    })
  }))

  validation {
    condition     = alltrue([for k, v in var.private_endpoints : strcontains(k, "pep-")])
    error_message = "Private endpoints should have the `pep-` prefix"
  }
}

variable "purge_protection_enabled" {
  default     = null
  description = <<-DESCRIPTION
  Is Purge Protection enabled for this Key Vault?

  This needs to be `true` if using the vault for disk encryption.

  Once Purge Protection has been Enabled it's not possible to Disable it. Deleting the
  Key Vault with Purge Protection Enabled will schedule the Key Vault to be deleted
  (which will happen by Azure in the configured number of days, currently 90 days - which
  will be configurable in Terraform in the future).
  DESCRIPTION
  type        = bool
}

variable "required_tags" {
  type = object({
    App         = string
    Environment = string
    GBU         = string
    ITSM        = optional(string, "MANAGEMENT")
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
  description = "Azure Resource Group Name"
}

variable "sku" {
  default     = "standard"
  description = "The Name of the SKU used for this Key Vault. Possible values are standard and premium."
  type        = string
}

variable "soft_delete_retention_days" {
  type        = number
  description = <<-DESCRIPTION
  The number of days that items should be retained for once soft-deleted. This value can
  be between 7 and 90 days. This field can only be configured one time and cannot be
  updated.
  DESCRIPTION
  default     = 90
}

variable "public_network_access_enabled" {
  default     = false
  description = "Whether public network access is allowed for this Key Vault."
  type        = bool
}
