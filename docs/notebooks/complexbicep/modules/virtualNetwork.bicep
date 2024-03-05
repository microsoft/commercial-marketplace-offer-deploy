param location string
param vNetName string
param vNetAddressPrefixes string
param vNetSubnetName string
param vNetSubnetAddressPrefix string
param networkSecurityGroupVnet string

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2020-05-01' = {
  name: vNetName
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        vNetAddressPrefixes
      ]
    }
    subnets: [
      {
        name: vNetSubnetName
        properties: {
          addressPrefix: vNetSubnetAddressPrefix
          networkSecurityGroup: {
            id: resourceId('Microsoft.Network/networkSecurityGroups', networkSecurityGroupVnet)
          }
        }
      }
    ]
  }
}

output vnetSubnetId string = '${virtualNetwork.id}/subnets/${vNetSubnetName}'

