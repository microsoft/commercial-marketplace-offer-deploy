
param unique string = uniqueString(utcNow())


var dependsOnName = format('{0}{1}', 'storedep0', substring(unique, 0, 5))

resource storageAccountWithDependsOn 'Microsoft.Storage/storageAccounts@2022-09-01' = {
  name: dependsOnName
  location: resourceGroup().location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'Storage'
}
