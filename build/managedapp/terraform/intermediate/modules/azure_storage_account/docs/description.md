# azure_storage

This module allows you to create a storage account with containers and private endpoints.
You can also control the storage account's firewall. In most cases, public network
access should be disabled and a private endpoint should be used.

## Private Endpoint

The VNet/subnet that the private endpoint utilizes must be in the same subscription as
the storage account. The private endpoint's DNS integration will be handled
automatically by an Azure Policy.

## Additional Info

* [Storage account overview](https://learn.microsoft.com/en-us/azure/storage/common/storage-account-overview)
* [Introduction to Azure Blob Storage](https://learn.microsoft.com/en-us/azure/storage/blobs/storage-blobs-introduction)
* [azurerm_storage_account](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account)
* [azurerm_storage_container](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_container)
* [azurerm_private_endpoint](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_endpoint)
