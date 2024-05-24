# azure_linux_web_app

This module allows you to create a Linux web app with a private endpoint and a custom
domain.

> NOTE: This module supports a very limited set of functionality. See the resource's
> documentation for a full feature list.

## Custom domain names & certificates

Before you add a custom domain name, you must create a CNAME for the custom domain that
points to the Azure web app. For example, if you want to add the custom domain
`mikentest.parsons.com`, you need to create the CNAME record
`mikentest.parsons.com` with the value `app-mike-test-1.azurewebsites.net`. If you don't
add this beforehand, the Terraform deployment will fail.

Each custom domain name will automatically have an App Service managed certificate
created for it. The type is automatically set to `SNI SSL`. It's possible to use your
own certificate with a custom domain, but this module doesn't support that.

The cert is automatically renewed continuously in six-month increments, 45 days before
expiration, as long as the prerequisites that you set up stay the same. All the
associated bindings are updated with the renewed certificate. [Source](
  https://learn.microsoft.com/en-us/azure/app-service/configure-ssl-certificate?tabs=apex#create-a-free-managed-certificate
)

## Private Endpoint

The VNet/subnet that the private endpoint utilizes must be in the same subscription as
the Linux web app. The private endpoint's DNS integration will be handled
automatically by an Azure Policy.

## Additional Info

* [azurerm_linux_web_app](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_web_app)
* [azurerm_private_endpoint](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_endpoint)
* [Using Private Endpoints for App Service apps](https://learn.microsoft.com/en-us/azure/app-service/overview-private-endpoint)
* [azurerm_app_service_custom_hostname_binding](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/app_service_custom_hostname_binding)
* [azurerm_app_service_managed_certificate](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/app_service_managed_certificate)
* [azurerm_app_service_certificate_binding](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/app_service_certificate_binding)
* [Map an existing custom DNS name to Azure App Service](https://learn.microsoft.com/en-us/azure/app-service/app-service-web-tutorial-custom-domain)
* [Secure a custom DNS name with a TLS/SSL binding in Azure App Service](https://learn.microsoft.com/en-us/azure/app-service/configure-ssl-bindings)
