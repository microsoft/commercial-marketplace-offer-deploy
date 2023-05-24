param name1 string = 'storetest09sdkjf34'
param name2 string = 'bobjacresource12'

var location = resourceGroup().location
var varName2 = name2

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
    storageAccountName:  varName2
    location: location
  }
  dependsOn: [
    mod1
  ]
}
