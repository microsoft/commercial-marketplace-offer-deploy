variable "location" {
  description = "The Azure region where the resource should be created"
  type        = string
}

variable "maximum_elastic_worker_count" {
  default     = null
  description = <<-DESCRIPTION
  The maximum number of workers to use in an Elastic SKU Plan. Cannot be set unless
  using an Elastic SKU.
  DESCRIPTION
  type        = number
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
  description = "The name of the resource. Changing this forces a new resource to be created."
  type        = string

  validation {
    condition     = strcontains(var.name, "asp-")
    error_message = "The name of an App Service Plan should start with the `asp-` prefix"
  }
}

variable "os_type" {
  description = <<-DESCRIPTION
  The O/S type for the App Services to be hosted in this plan. Possible values include
  Windows, Linux, and WindowsContainer. Changing this forces a new resource to be created.
  DESCRIPTION
  type        = string
}

variable "per_site_scaling_enabled" {
  default     = false
  description = "Should Per Site Scaling be enabled."
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
  The SKU for the plan. Possible values include `B1`, `B2`, `B3`, `D1`, `F1`, `I1`,
  `I2`, `I3`, `I1v2`, `I2v2`, `I3v2`, `I4v2`, `I5v2`, `I6v2`, `P1v2`, `P2v2`, `P3v2`,
  `P0v3`, `P1v3`, `P2v3`, `P3v3`, `P1mv3`, `P2mv3`, `P3mv3`, `P4mv3`, `P5mv3`, `S1`,
  `S2`, `S3`, `SHARED`, `EP1`, `EP2`, `EP3`, `WS1`, `WS2`, `WS3`, and `Y1`.

  Isolated SKUs (`I1`, `I2`, `I3`, `I1v2`, `I2v2`, and `I3v2`) can only be used with
  App Service Environments.
  Elastic and Consumption SKUs (`Y1`, `EP1`, `EP2`, and `EP3`) are for use with
  Function Apps.
  DESCRIPTION
  type        = string
}

variable "worker_count" {
  default     = null
  description = "The number of Workers (instances) to be allocated."
  type        = number
}

variable "zone_balancing_enabled" {
  default     = null
  description = <<-DESCRIPTION
  Should the Service Plan balance across Availability Zones in the region. Changing this
  forces a new resource to be created.

  If this setting is set to `true` and the `worker_count` value is specified, it should
  be set to a multiple of the number of availability zones in the region. Please see the
  Azure documentation for the number of Availability Zones in your region.
  DESCRIPTION
  type        = bool
}
