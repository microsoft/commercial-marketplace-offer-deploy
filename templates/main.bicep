param location string = resourceGroup().location
param serviceBusNamespaceName string = 'modm-servicebus'

module servicebusModule 'modules/servicebus.bicep' = {
  name: 'servicebusDeploy'
  params: {
    location: location
    serviceBusNamespaceName: serviceBusNamespaceName
  }
}
