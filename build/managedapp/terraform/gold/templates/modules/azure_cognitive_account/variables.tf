variable "custom_subdomain_name" {
  description = <<-DESCRIPTION
  The subdomain name used for token-based authentication. This property is required when
  a private endpoint is used. Changing this forces a new resource to be created.
  DESCRIPTION
  type        = string
}

variable "dynamic_throttling_enabled" {
  default     = null
  description = "Whether to enable the dynamic throttling for this Cognitive Service Account."
  type        = bool
}

variable "fqdns" {
  default     = null
  description = "List of FQDNs allowed for the Cognitive Account."
  type        = list(string)
}

variable "kind" {
  description = <<-DESCRIPTION
  Specifies the type of Cognitive Service Account that should be created. For a list of
  possible values see:
  https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/cognitive_account#kind
  DESCRIPTION
  type        = string
}

variable "local_auth_enabled" {
  default     = true
  description = "Whether local authentication methods is enabled for the Cognitive Account."
  type        = bool
}

variable "location" {
  description = "The Azure region where the resource should be created"
  type        = string
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

variable "name" {
  description = "The name of the Cognitive Services account"
  type        = string
}

variable "outbound_network_access_restricted" {
  default     = false
  description = "Whether outbound network access is restricted for the Cognitive Account."
  type        = bool
}

variable "private_endpoints" {
  default     = {}
  description = <<-DESCRIPTION
  A collection of private endpoints to create & connect to the cognitive account.
  `subresource_names` defaults to `['account']`, as that is the only available
  subresource. The VNet/subnet that is looked up must be in the same subscription as the
  cognitive account.
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

variable "public_network_access_enabled" {
  default     = false
  description = "Whether public network access is allowed for the Cognitive Account."
  type        = bool
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

variable "resource_group_name" {
  description = "The resource group in which the resource should be created"
  type        = string
}

variable "sku_name" {
  description = <<-DESCRIPTION
  Specifies the SKU Name for this Cognitive Service Account. For possible values see:
  https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/cognitive_account#sku_name
  DESCRIPTION
  type        = string
}
