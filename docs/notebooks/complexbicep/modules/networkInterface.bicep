param location string
param networkInterfaceName string
param publicIPAddressName string
param vNetName string
param vNetSubnetName string

resource networkInterface 'Microsoft.Network/networkInterfaces@2020-05-01' = {
  name: networkInterfaceName
  location: location
  properties: {
    ipConfigurations: [
      {
        name: 'ipconfig1'
        properties: {
          privateIPAllocationMethod: 'Dynamic'
          publicIPAddress: {
            id: resourceId('Microsoft.Network/publicIPAddresses', publicIPAddressName)
          }
          subnet: {
            id: resourceId('Microsoft.Network/virtualNetworks/subnets', vNetName, vNetSubnetName)
          }
        }
      }
    ]
  }
}

output networkInterfaceId string = networkInterface.id
