# azure_app_service_plan

This module allows you to create an Azure App Service plan

> NOTE: This module doesn't currently support placing the Service Plan in an
> [App Service Environment](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/service_plan#app_service_environment_id)

## Network Integration

It's possible to set up a VNet integration for an App Service Plan, but this module
doesn't support that because a private endpoint should be configured per web app.

> Virtual network integration gives your app access to resources in your virtual
> network, but it doesn't grant inbound private access to your app from the virtual
> network. Private site access refers to making an app accessible only from a private
> network, such as from within an Azure virtual network. Virtual network integration is
> used only to make outbound calls from your app into your virtual network. Refer to
> private endpoint for inbound private access.

## Additional Info

* [azurerm_service_plan](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/service_plan)
* [Integrate your app with an Azure virtual network](https://learn.microsoft.com/en-us/azure/app-service/overview-vnet-integration)
