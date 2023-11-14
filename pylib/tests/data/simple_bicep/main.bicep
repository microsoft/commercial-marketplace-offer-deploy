param location string

module stgModule 'modules/storageAccount.bicep' = {
  name: 'storageDeploy'
  params: {
    location: location
  }
} 

output storageAccountId string = stgModule.outputs.storageAccountId
