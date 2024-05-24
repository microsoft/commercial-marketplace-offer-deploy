variable "authentication_failure_mode" {
  default     = null
  description = <<-DESCRIPTION
  Specifies the response that the Search Service should return for requests that fail
  authentication. Possible values include `http401WithBearerChallenge` or `http403`.

  This can only be configured when `local_authentication_enabled` is set to true - which
  when set together specifies that both API Keys and AzureAD Authentication should be
  supported.
  DESCRIPTION
  type        = string
}

variable "hosting_mode" {
  default     = "default"
  description = <<-DESCRIPTION
  Specifies the Hosting Mode, which allows for High Density partitions (that allow for
  up to 1000 indexes) should be supported. Possible values are highDensity or default.
  Changing this forces a new Search Service to be created.

  This can only be configured when `sku` is set to `standard3`.
  DESCRIPTION
  type        = string
}

variable "local_authentication_enabled" {
  default     = true
  description = "Specifies whether the Search Service allows authenticating using API Keys"
  type        = bool
}

variable "location" {
  description = "The Azure region where the resource should be created"
  type        = string
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
  description = <<-DESCRIPTION
  The name of the Search Service. Changing this forces a new resource to be created.
  DESCRIPTION
  type        = string

  validation {
    condition     = strcontains(var.name, "srch-")
    error_message = "The name of a Search Service should start with the `srch-` prefix"
  }
}

variable "partition_count" {
  default     = 1
  description = <<-DESCRIPTION
  Specifies the number of partitions which should be created. This field cannot be set
  when using a free or basic sku (see the Microsoft documentation below). Possible values
  include 1, 2, 3, 4, 6, or 12.
  https://learn.microsoft.com/azure/search/search-sku-tier
  DESCRIPTION
  type        = number
}

variable "private_endpoints" {
  default     = {}
  description = <<-DESCRIPTION
  A collection of private endpoints to create & connect to the search service.
  `subresource_names` defaults to `['searchService']`, as that is the only available
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
  description = "Whether public network access is allowed for this resource."
  type        = bool
}

variable "replica_count" {
  default     = null
  description = <<-DESCRIPTION
  Specifies the number of Replica's which should be created for this Search Service.
  This field cannot be set when using a free sku (see the Microsoft documentation).
  https://learn.microsoft.com/azure/search/search-sku-tier
  DESCRIPTION
  type        = number
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

variable "semantic_search_sku" {
  default     = null
  description = <<-DESCRIPTION
  Specifies the Semantic Search SKU which should be used for this Search Service.
  Possible values include `free` and `standard`. This cannot be defined if your Search
  Services `sku` is set to `free`. The Semantic Search feature is only available in
  certain regions, please see the product documentation for more information:
  https://learn.microsoft.com/azure/search/semantic-search-overview#availability-and-pricing
  DESCRIPTION
  type        = string
}

variable "sku" {
  description = <<-DESCRIPTION
  The SKU which should be used for this Search Service. Possible values include `basic`,
  `free`, `standard`, `standard2`, `standard3`, `storage_optimized_l1` and
  `storage_optimized_l2`. Changing this forces a new Search Service to be created.

  The basic and free SKUs provision the Search Service in a Shared Cluster - the
  standard SKUs use a Dedicated Cluster.

  The SKUs `standard2`, `standard3`, `storage_optimized_l1` and `storage_optimized_l2`
  are only available by submitting a quota increase request to Microsoft. Please see the
  product documentation on how to submit a quota increase request:
  https://learn.microsoft.com/azure/azure-resource-manager/troubleshooting/error-resource-quota?tabs=azure-cli
  DESCRIPTION
  type        = string
}
