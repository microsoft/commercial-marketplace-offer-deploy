param location string = resourceGroup().location

param appVersion string = 'v0.1.8'

@description('admin email used for Lets Encrypt.')
param adminEmail string

var containerImage = 'ghcr.io/gpsuscodewith/modm:${appVersion}'

module servicebusModule 'modules/servicebus.bicep' = {
  name: 'serviceBus'
  params: {
    location: location
    appVersion: appVersion
  }
}

module containerInstanceModule 'modules/containerInstance.bicep' = {
  name: 'containerInstance'
  params: {
    location: location
    appVersion: appVersion
    containerImage: containerImage
    resourceGroupName: resourceGroup().name
    subscriptionId: subscription().subscriptionId
    tenantId: subscription().tenantId
    acmeEmail: adminEmail
    serviceBusNamespace: servicebusModule.outputs.serviceBusNamespace
  }
  dependsOn: [
    servicebusModule
  ]
}

module roleAssignments 'modules/roleAssignments.bicep' = {
  name: 'roleAssignments'
  params: {
    containerGroupName: containerInstanceModule.outputs.containerGroupName
    serviceBusNamespace: servicebusModule.outputs.serviceBusNamespace
    storageAccountName: containerInstanceModule.outputs.storageAccountName
  }
  dependsOn: [
    servicebusModule
    containerInstanceModule
  ]
}
