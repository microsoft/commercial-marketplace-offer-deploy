
param location string = 'eastus'
param suffix string = substring(uniqueString(utcNow()), 0, 5)

param storageCount int = 5

resource storageAccounts 'Microsoft.Storage/storageAccounts@2022-09-01' = [for i in range(0, storageCount): {
  name: '${i}storage${uniqueString(resourceGroup().id)}'
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'Storage'
}]

var dependsOnName = format('{0}{1}', 'storageAccounts-', suffix)

resource storageAccountWithDependsOn 'Microsoft.Storage/storageAccounts@2022-09-01' = {
  name: dependsOnName
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'Storage'
  dependsOn: [
    storageAccounts
  ]
}
