output "id" {
  description = "The ID of the Service Plan."
  value       = azurerm_service_plan.this.id
}

output "reserved" {
  description = <<-DESCRIPTION
  Whether this is a reserved Service Plan Type.
 `true` if `os_type` is `Linux`, otherwise `false`.
  DESCRIPTION
  value       = azurerm_service_plan.this.reserved
}
