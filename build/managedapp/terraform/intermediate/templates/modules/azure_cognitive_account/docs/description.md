# azure_cognitive_account

This module allows you to deploy a cognitive account with a private endpoint.

> NOTE: You must create your first Face, Text Analytics, or Computer Vision resources
> from the Azure portal to review and acknowledge the terms and conditions. In Azure
> Portal, the checkbox to accept terms and conditions is only displayed when a US region
> is selected. More information on [Prerequisites](
  https://docs.microsoft.com/azure/cognitive-services/cognitive-services-apis-create-account-cli?tabs=windows#prerequisites
).

## Private Endpoint

The custom subdomain name is required due to the usage of private endpoints.
The VNet/subnet that the private endpoint utilizes must be in the same subscription as
the cognitive services account. The private endpoint's DNS integration will be handled
automatically by an Azure Policy.

## Additional Info

* [azurerm_cognitive_account](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/cognitive_account)
* [azurerm_private_endpoint](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_endpoint)
