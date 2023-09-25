output "networkID_out" {
    value = azurerm_network_interface.interface_name.id
}

output "networkName_out" {
    value = azurerm_virtual_network.vNet.name
}

output "networkVNetID_out" {
    value = azurerm_virtual_network.vNet.id
}


