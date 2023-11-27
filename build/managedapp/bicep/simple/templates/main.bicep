param location string
param resourceGroupName string

module stgModule 'modules/storageAccount.bicep' = {
  name: 'storageDeploy'
  params: {
    location: location
  }
} 

output storageAccountId string = stgModule.outputs.storageAccountId
