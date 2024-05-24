provider "azurerm" {
  environment                = var.environment
  skip_provider_registration = true
  subscription_id            = var.subscription_id

  features {}
}
