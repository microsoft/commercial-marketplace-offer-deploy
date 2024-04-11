param location string
param resourceGroupName string
param artifactsLocationSasToken string

module stgModule 'modules/storageAccount.bicep' = {
  name: 'storageDeploy'
  params: {
    location: location
  }
} 

output storageAccountId string = stgModule.outputs.storageAccountId
