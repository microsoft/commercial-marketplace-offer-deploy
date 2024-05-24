variable "admin_username" {
  description = "The admin username for the PostgreSQL flexible server"
  type        = string
}

variable "auto_grow_enabled" {
  default     = false
  description = "Is the storage auto grow for PostgreSQL Flexible Server enabled?"
  type        = bool
}

variable "backup_retention_days" {
  default     = 7
  description = <<-DESCRIPTION
  The backup retention days for the PostgreSQL Flexible Server. Possible values are
  between 7 and 35 days.
  DESCRIPTION
  type        = number
}

variable "geo_redundant_backup_enabled" {
  default     = false
  description = <<-DESCRIPTION
  Is Geo-Redundant backup enabled on the PostgreSQL Flexible Server. Changing this
  forces a new PostgreSQL Flexible Server to be created.
  DESCRIPTION
  type        = bool
}

variable "high_availability" {
  default     = null
  description = <<-DESCRIPTION
  `mode` - The high availability mode for the PostgreSQL Flexible Server. Possible values
  are SameZone or ZoneRedundant.
  `standby_availability_zone` - Specifies the Availability Zone in which the standby
  Flexible Server should be located.
  DESCRIPTION
  type = object({
    mode                      = string
    standby_availability_zone = optional(number)
  })
}

variable "is_gov" {
  description = "Is the deployment running within the Azure US Government tenant"
  type        = bool
}

variable "location" {
  description = "The Azure region where the resource should be created"
  type        = string
}

variable "maintenance_window" {
  default     = null
  description = <<-DESCRIPTION
  `day_of_week` - The day of week for the maintenance window, where the week starts on a
  Sunday, i.e. Sunday = 0, Monday = 1.
  `start_hour` - The start hour for the maintenance window
  `start_minute` - The start minute for the maintenance window
  DESCRIPTION
  type = object({
    day_of_week  = optional(number, 0)
    start_hour   = optional(number, 0)
    start_minute = optional(number, 0)
  })
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
  The name which should be used for this PostgreSQL Flexible Server. Changing this forces
  a new PostgreSQL Flexible Server to be created. his must be unique across the entire
  Azure service, not just within the resource group.
  DESCRIPTION
  type        = string
}

variable "postgresql_version" {
  description = <<-DESCRIPTION
  The version of PostgreSQL Flexible Server to use. Possible values are 11,12, 13, 14,
  15 and 16.
  DESCRIPTION
  type        = string
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
  The SKU Name for the PostgreSQL Flexible Server. The name of the SKU follows the tier +
  name pattern (e.g. B_Standard_B1ms, GP_Standard_D2s_v3, MO_Standard_E4s_v3). See:
  https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-compute-storage
  DESCRIPTION
  type        = string
}

variable "storage_mb" {
  default     = 32768
  description = <<-DESCRIPTION
  The max storage allowed for the PostgreSQL Flexible Server. This can only be scaled up.
  Possible values are 32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4193280,
  4194304, 8388608, 16777216 and 33553408.
  DESCRIPTION
  type        = number
}

variable "storage_tier" {
  default     = null
  description = <<-DESCRIPTION
  The name of storage performance tier for IOPS of the PostgreSQL Flexible Server.
  Possible values are P4, P6, P10, P15,P20, P30,P40, P50,P60, P70 or P80. Default value
  is dependant on the storage_mb value. See:
  https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-compute-storage
  DESCRIPTION
  type        = string
}

variable "vnet_integration" {
  default     = null
  description = <<-DESCRIPTION
  This object must be set if you wish to set up the virtual network integration. The
  VNet, subnet, and their resource group must be in the same subscription as the
  PostgreSQL Flexible Server. The subnet should not have any other resource deployed in
  it and must be delegated to the PostgreSQL Flexible Server. Changing this forces a new
  PostgreSQL Flexible Server to be created.
  DESCRIPTION
  type = object({
    delegated_subnet_name    = string
    vnet_name                = string
    vnet_resource_group_name = string
  })
}

variable "zone" {
  default     = null
  description = <<-DESCRIPTION
  Specifies the Availability Zone in which the PostgreSQL Flexible Server should be
  located. Azure will automatically assign an Availability Zone if one is not specified.
  If the PostgreSQL Flexible Server fails-over to the Standby Availability Zone, the zone
  will be updated to reflect the current Primary Availability Zone.

  Terraform has been configured to ignore changes to the `zone` and
  `high_availability.0.standby_availability_zone` fields, which means it will not migrate
  the PostgreSQL Flexible Server back to it's primary Availability Zone after a fail-over.
  DESCRIPTION
  type        = number
}
