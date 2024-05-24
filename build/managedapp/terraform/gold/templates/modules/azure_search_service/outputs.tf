output "id" {
  description = "The ID of the Search Service."
  value       = azurerm_search_service.this.id
}

output "primary_key" {
  description = "The Primary Key used for Search Service Administration."
  value       = azurerm_search_service.this.primary_key
}

output "query_keys" {
  description = "The Search Service's query keys."
  value       = azurerm_search_service.this.query_keys
}

output "secondary_key" {
  description = "The Secondary Key used for Search Service Administration."
  value       = azurerm_search_service.this.secondary_key
}
