param name string = ''
param location string = ''

// Storage account name must be between 3 and 24 characters in length and use numbers and lower-case letters only.
//var storageAccountName = take('st${name}${location}', 24)
var storageAccountName = name

module storageAccount 'storage/storageAccounts.bicep' = {
  name: 'storageAccountDeploy'
  params: {
    name: storageAccountName
    location: location
  }
}

module container 'storage/storageAccounts/blobServices/containers.bicep' = {
  name: 'containerDeploy'
  params: {
    name: 'pulp'
    parent: storageAccountName
  }
  dependsOn: [
    storageAccount
  ]
}


output name string = storageAccountName
