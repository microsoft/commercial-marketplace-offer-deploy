@description('Name of the Service Bus namespace')
param appVersion string

var serviceBusQueueNames = [
  'events'
  'operations'
  'healthcheck'
]

@description('Location for all resources.')
param location string = resourceGroup().location

var versionSuffix = replace(appVersion, '.', '')
var serviceBusNamespace = 'modmsb-${versionSuffix}'

resource serviceBus 'Microsoft.ServiceBus/namespaces@2022-01-01-preview' = {
  name: serviceBusNamespace
  location: location
  sku: {
    name: 'Standard'
  }
  properties: {}
}

resource serviceBusQueues 'Microsoft.ServiceBus/namespaces/queues@2022-01-01-preview' = [for i in range(0, length(serviceBusQueueNames)): {
  parent: serviceBus
  name: serviceBusQueueNames[i]
  properties: {
    lockDuration: 'PT5M'
    maxSizeInMegabytes: 1024
    requiresDuplicateDetection: false
    requiresSession: false
    defaultMessageTimeToLive: 'P10675199DT2H48M5.4775807S'
    deadLetteringOnMessageExpiration: false
    duplicateDetectionHistoryTimeWindow: 'PT10M'
    maxDeliveryCount: 10
    autoDeleteOnIdle: 'P10675199DT2H48M5.4775807S'
    enablePartitioning: false
    enableExpress: false
  }
}]

output serviceBusNamespace string = serviceBusNamespace
