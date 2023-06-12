
param location string = 'eastus'
param unique string = uniqueString(utcNow())

var storageCount = 5

resource storageAccounts 'Microsoft.Storage/storageAccounts@2022-09-01' = [for i in range(0, storageCount): {
  name: '${i}stor0${substring(unique, i, 5)}'
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'Storage'
}]
