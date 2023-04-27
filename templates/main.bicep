param location string = resourceGroup().location

param appVersion string = 'latest'

@description('admin email used for Lets Encrypt.')
param acmeEmail string

module servicebusModule './modules/servicebus.bicep' = {
  name: 'serviceBus'
  params: {
    location: location
    appVersion: appVersion
  }
}

var containerImage = 'ghcr.io/gpsuscodewith/modm'

var roleIds = {
  reader: 'acdd72a7-3385-48ef-bd42-f606fba81ae7'
  owner: '8e3af657-a8ff-443c-a75c-2fe8c4bcb635'
  storageAccountContributor: '0c867c2a-1d8c-454a-a3db-ab2ea1bdc8bb'

  serviceBusDataReceiver: '4f6d3b9b-027b-4f4c-9142-0e5a2a2247e0'
  serviceBusDataSender: '69a216fc-b8fb-44d8-bc22-1f3c2cd27a39'
  serviceBusDataOwner: '090c5cfd-751d-490a-894a-3ce6f1109419'
  
  eventGridContributor: '1e241071-0855-49ea-94dc-649edcd759de'
  eventGridDataSender: 'd5a91429-5739-47e2-a06b-3470a27159e7'
  eventGridEventSubscriptionContributor: '428e0ff0-5e57-4d9c-a221-2c70d0e0a443'
  eventGridEventSubscriptionReader: '2414bbcf-6497-4faf-8c65-045460748405'
}

var roleAssignmentIds = {
  owner: guid(resourceGroup().id, appVersion, roleIds.owner)
  reader: guid(resourceGroup().id, appVersion, roleIds.reader)
}


module containerInstanceModule './modules/containerInstance.bicep' = {
  name: 'containerInstance'
  params: {
    location: location
    appVersion: appVersion
    containerImage: containerImage
    resourceGroupName: resourceGroup().name
    subscriptionId: subscription().subscriptionId
    tenantId: subscription().tenantId
    acmeEmail: acmeEmail
    serviceBusNamespace: servicebusModule.outputs.serviceBusNamespace
    roleAssignmentIds: roleAssignmentIds
  }
  dependsOn: [
    servicebusModule
  ]
}

module roleAssignments './modules/roleAssignments.bicep' = {
  name: 'roleAssignments'
  params: {
    containerGroupName: containerInstanceModule.outputs.containerGroupName
    serviceBusNamespace: servicebusModule.outputs.serviceBusNamespace
    storageAccountName: containerInstanceModule.outputs.storageAccountName
    roleAssignmentIds: roleAssignmentIds
  }
  dependsOn: [
    servicebusModule
    containerInstanceModule
  ]
}
