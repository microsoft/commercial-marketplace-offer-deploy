variable "name" {
  description = "The name of the zone. Must be a valid domain name"
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

variable "virtual_network_links" {
  default     = {}
  description = <<-DESCRIPTION
  A collection of virtual network links to attach the private DNS zone to one or more
  virtual networks. This will allow a VNet to resolve records in this zone. They key
  should be the name of the virtual network and it must be in the same subscription as
  the private DNS zone. `auto_registration_enabled` allows virtual machines in the VNet
  to automatically have DNS records created in the zone.
  DESCRIPTION
  type = map(object({
    auto_registration_enabled = optional(bool, false)
    resource_group_name       = string
  }))
}
