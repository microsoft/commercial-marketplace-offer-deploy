# azure_search_service

This module allows you to deploy an Azure AI Search Service with a private endpoint.

> NOTE: Semantic Ranking isn't supported in Azure Government yet
>
> NOTE: This module doesn't currently support setting a Search Service's
> [identity](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/search_service#identity)

## Private Endpoint

The VNet/subnet that the private endpoint utilizes must be in the same subscription as
the cognitive services account. The private endpoint's DNS integration will be handled
automatically by an Azure Policy.

## Additional Info

* [azurerm_search_service](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/search_service)
* [azurerm_private_endpoint](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_endpoint)
* [Create a Private Endpoint for a secure connection to Azure AI Search](https://learn.microsoft.com/en-us/azure/search/service-create-private-endpoint#create-a-search-service-with-a-private-endpoint)
* [Semantic ranking in Azure AI Search](https://learn.microsoft.com/en-us/azure/search/semantic-search-overview)
