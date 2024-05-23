variable "custom_domains" {
  default     = []
  description = <<-DESCRIPTION
  A set of custom domains to assign to the web app. Each custom domain will automatically
  have an App Service managed certificate created and bound to the custom hostname.

  > NOTE: Before you attempt to add custom domains, you must add a CNAME for the custom
  > domain which points to the web app using its Azure hostname.
  DESCRIPTION
  type        = set(string)
}

variable "enabled" {
  default     = true
  description = "Should the Linux Web App be enabled?"
  type        = bool
}

variable "https_only" {
  default     = false
  description = "Should the Linux Web App require HTTPS connections."
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
  The name which should be used for this Linux Web App. Changing this forces a new Linux
  Web App to be created.

  Terraform will perform a name availability check as part of the creation process. If
  this Web App is part of an App Service Environment Terraform will require Read
  permission on the ASE for this to complete reliably.
  DESCRIPTION
  type        = string

  validation {
    condition     = strcontains(var.name, "app-")
    error_message = "The name of web app should start with the `app-` prefix"
  }
}

variable "private_endpoints" {
  default     = {}
  description = <<-DESCRIPTION
  A collection of private endpoints to create & connect to the Linux web app.
  `subresource_names` defaults to `['sites']`, as that is the only available
  subresource. The VNet/subnet that is looked up must be in the same subscription as the
  Linux web app.
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
  description = "Should public network access be enabled for the Web App."
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

variable "service_plan_id" {
  description = "The ID of the Service Plan that this Linux App Service will be created in."
  type        = string
}

variable "site_config" {
  description = <<-DESCRIPTION
  `always_on` - If this Linux Web App is Always On enabled. This must be explicitly set
    to `false` when using `Free`, `F1`, `D1`, or `Shared` service plans.
  `application_stack`:
    Only one runtime can be specified.

    `dotnet_version` - The version of .NET to use. Possible values include `3.1`, `5.0`
      `6.0`, `7.0` and `8.0`.
    `go_version` - The version of Go to use. Possible values include `1.18`, and `1.19`.
    `java_server` - The Java server type. Possible values include `JAVA`, `TOMCAT`, and `JBOSSEAP`.
    `java_server_version` - The Version of the java_server to use.
    `java_version` - The Version of Java to use. Possible values include `8`, `11`, and `17`.

    The valid version combinations for `java_version`, `java_server` and
    `java_server_version` can be checked from the command line via `az webapp list-runtimes --linux`.

    `node_version` - The version of Node to run. Possible values include `12-lts`, `14-lts`,
      `16-lts`, `18-lts` and `20-lts`. This property conflicts with java_version.
    `php_version` - The version of PHP to run. Possible values are `7.4`, `8.0`, `8.1` and `8.2`.
    `python_version` - The version of Python to run. Possible values include `3.7`, `3.8`,
      `3.9`, `3.10`, `3.11` and `3.12`.
    `ruby_version` - The version of Ruby to run. Possible values include `2.6` and `2.7`.
  `ftps_state` - The State of FTP / FTPS service. Possible values include `AllAllowed`,
    `FtpsOnly`, and `Disabled`.
  `http2_enabled` - Should the HTTP2 be enabled?
  `ip_restriction_default_action` -   The Default action for traffic that does not match
    any ip_restriction rule. Possible values include Allow and Deny.
  `load_balancing_mode` - The Site load balancing. Possible values include:
    `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`,
    `RequestHash`, `PerSiteRoundRobin`.
  `minimum_tls_version` - The configures the minimum version of TLS required for SSL
    requests. Possible values include: `1.0`, `1.1`, and `1.2`.
  `scm_ip_restriction_default_action` -   The Default action for traffic that does not
    match any scm_ip_restriction rule. Possible values include Allow and Deny.
  `use_32_bit_worker` - Should the Linux Web App use a 32-bit worker?
  `vnet_route_all_enabled` - Should all outbound traffic have NAT Gateways, Network
    Security Groups and User Defined Routes applied?
  DESCRIPTION
  type = object({
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
}

variable "webdeploy_publish_basic_authentication_enabled" {
  default     = true
  description = <<-DESCRIPTION
  Should the default WebDeploy Basic Authentication publishing credentials enabled.

  Setting this value to `true` will disable the ability to use `zip_deploy_file` which
  currently relies on the default publishing profile.
  DESCRIPTION
  type        = bool
}
