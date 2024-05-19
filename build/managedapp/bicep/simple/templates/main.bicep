param location string
param resourceGroupName string
param tier string
param artifactsLocationSasToken string

module stgModule 'modules/storageAccount.bicep' = {
  name: 'storageDeploy'
  params: {
    location: location
    tier: tier
  }
} 

output storageAccountId string = stgModule.outputs.storageAccountId
