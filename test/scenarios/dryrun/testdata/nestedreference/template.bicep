param name1 string = 'Standard_LRS'
param name2 string = 'Standard_LRS'

var location = resourceGroup().location

module mod1 './modules/storageAccount.bicep' = {
  name: 'mod1'
  params: {
    storageAccountName: name1
    storageAccountName2: name2
    location: location
  }
}

module mod2 './modules/storageAccountWithRefence.bicep' = {
  name: 'mod2'
  params: {
    storageAccountName: mod1.outputs.storageAccountName2
    location: location
  }
  dependsOn: [
    mod1
  ]
}
