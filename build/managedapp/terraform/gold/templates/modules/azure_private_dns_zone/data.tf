data "azurerm_virtual_network" "this" {
  for_each = var.virtual_network_links

  name                = each.key
  resource_group_name = each.value.resource_group_name
}
