output "endpoint" {
  description = "The endpoint used to connect to the cognitive services account"
  value       = azurerm_cognitive_account.this.id
}

output "id" {
  description = "The ID of the cognitive services account"
  value       = azurerm_cognitive_account.this.id
}
