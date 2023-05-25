
var firstName = 'bobjacresource1'
var firstName2 = '${firstName}2'
var location = 'eastus2'

module storage1 'modules/storagesmall.bicep' = {
  name: 'storage1'
  params: {
    name: firstName
    location: location
  }
}

module storage2 'modules/storagesmall.bicep' = {
  name: 'storage2'
  params: {
    name: firstName2
    location: location
  }
  dependsOn: [
    storage1
  ]
}


