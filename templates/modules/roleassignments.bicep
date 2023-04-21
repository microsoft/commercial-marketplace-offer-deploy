param storageAccountName string
param containerGroupName string
param serviceBusNamespace string
param appInsightsName string

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' existing = {
  name: containerGroupName
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' existing = {
  name: storageAccountName
}

resource serviceBus 'Microsoft.ServiceBus/namespaces@2022-01-01-preview' existing = {
  name: serviceBusNamespace
}

resource appInsights 'Microsoft.Insights/components@2020-02-02' existing = {
  name: appInsightsName
}

var roles = {
  storageAccountContributor: '0c867c2a-1d8c-454a-a3db-ab2ea1bdc8bb'
  serviceBusDataReceiver: '4f6d3b9b-027b-4f4c-9142-0e5a2a2247e0'
  serviceBusDataSender: '69a216fc-b8fb-44d8-bc22-1f3c2cd27a39'
  appInsightsContributor: 'b24988ac-6180-42a0-ab88-20f7382dd24c'
}

resource roleAssignmentStorageAcct 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount //assigns to storage acct
  name: guid(storageAccount.id, containerGroup.name, roles.storageAccountContributor)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.storageAccountContributor)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource serviceBusReceiverAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: serviceBus
  name: guid(serviceBus.id, containerGroup.name, roles.serviceBusDataReceiver)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.serviceBusDataReceiver)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource serviceBusSenderAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: serviceBus
  name: guid(serviceBus.id, containerGroup.name, roles.serviceBusDataSender)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.serviceBusDataSender)
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource roleAssignmentApplicationInsights 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' {
  scope: appInsights
  name: guid(appInsights.id, containerGroup.name, roles.appInsightsContributor)
  properties: {
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', roles.appInsightsContributor)
    principalId: containerGroup.identity.principalId
    principalType: 'ApplicationInsights'
  }
}
