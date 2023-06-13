
param location string = 'eastus'
param unique string = substring(uniqueString(utcNow()), 0, 5)


module storageAcounts './modules/storageAccounts.bicep' = {
  name: 'storageAccounts'
  params:{
    location: resourceGroup().location
  }
}

module dependsOn './modules/deploymentScript.bicep' = {
  name: 'dependsOnStorageAccounts'
  dependsOn: [
    storageAcounts
  ]
}
