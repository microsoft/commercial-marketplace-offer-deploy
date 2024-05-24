resource "azurerm_private_dns_zone" "this" {
  name                = var.name
  resource_group_name = var.resource_group_name
  tags                = var.required_tags
}

resource "azurerm_private_dns_zone_virtual_network_link" "this" {
  for_each = var.virtual_network_links

  name                  = each.key
  private_dns_zone_name = azurerm_private_dns_zone.this.name
  registration_enabled  = each.value.auto_registration_enabled
  resource_group_name   = var.resource_group_name
  tags                  = var.required_tags
  virtual_network_id    = data.azurerm_virtual_network.this[each.key].id
}
