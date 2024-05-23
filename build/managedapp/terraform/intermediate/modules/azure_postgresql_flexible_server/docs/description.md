# azure_postgresql_flexible_server

This module allows you to deploy a PostgreSQL Flexible Server with VNet integration.
Password authentication is used and the password is randomly generated.

## What's not supported yet

* Within the `azurerm_postgresql_flexible_server` resource:
  * [identity](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server#identity)
  * [custom_managed_key](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server#customer_managed_key)
  * [Active Directory authentication](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server#authentication)
    * Only password authentication is currently supported
  * [create_mode](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server#create_mode)
  other than default
    * As a result, [source_server_id](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server#source_server_id)
    also isn't supported
* [azurerm_postgresql_flexible_server_configuration](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server_configuration)
* [azurerm_postgresql_flexible_server_database](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server_database)
* [azurerm_postgresql_flexible_server_firewall_rule](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server_firewall_rule)
  * We shouldn't be allowing public access

## Private Link

VNet integration is supported and requires a subnet delegated to
`Microsoft.DBforPostgreSQL/flexibleServers`. The VNet integration will be automatically
connected to a corporate IT private DNS zone for centralized DNS. Configuring VNet
integration will automatically disable public access. If VNet integration is not used,
public access will remain enabled, but by default is unaccessible unless you whitelist
traffic from a public IP.

Private endpoints for PostgreSQL flexible servers are currently in preview which is not
available in US Gov Virginia, so this module doesn't support it yet.

## Additional Info

* [azurerm_postgresql_flexible_server](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server)
* [Compute and storage options in Azure Database for PostgreSQL - Flexible Server](https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-compute-storage#compute-tiers-vcores-and-server-types)
* [Networking overview for Azure Database for PostgreSQL - Flexible Server with private access (VNET Integration)](https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-networking-private)
* [Azure Database for PostgreSQL - Flexible Server networking with Private Link - Preview](https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-networking-private-link)
