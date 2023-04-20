
param principalID string 
param storageAccountName string = ''
param containerGroupName string = ''


resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' existing = {
  name: storageAccountName
}

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' existing = {
  name: containerGroupName
}

resource storageAccountContributorRoleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' existing = {
  scope: subscription()
  name: '0c867c2a-1d8c-454a-a3db-ab2ea1bdc8bb'
}

resource serviceBusDataReceiverRoleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' existing = {
  scope: subscription()
  name: '4f6d3b9b-027b-4f4c-9142-0e5a2a2247e0'
}

resource serviceBusDataSenderRoleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' existing = {
  scope: subscription()
  name: '69a216fc-b8fb-44d8-bc22-1f3c2cd27a39'
}

resource applicationInsightsRoleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' existing {
  scope: subscription()
  name: 'b24988ac-6180-42a0-ab88-20f7382dd24c'
}

resource roleAssignmentStorageAcct 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount //assigns to storage acct
  name: guid(storageAccount.id, containerGroup.name, storageAccountContributorRoleDefinition.id)
  properties: {
    roleDefinitionId: storageAccountContributorRoleDefinition.id
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource roleAssignmentServiceBusReceiver 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount //assigns to storage acct
  name: guid(storageAccount.id, containerGroup.name, serviceBusDataReceiverRoleDefinition.id)
  properties: {
    roleDefinitionId: serviceBusDataReceiverRoleDefinition.id
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource roleAssignmentServiceBusSender 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount //assigns to storage acct
  name: guid(storageAccount.id, containerGroup.name, serviceBusDataSenderRoleDefinition.id)
  properties: {
    roleDefinitionId: serviceBusDataSenderRoleDefinition.id
    principalId: containerGroup.identity.principalId
    principalType: 'ServicePrincipal'
  }
}

resource roleAssignmentApplicationInsights 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' {
  scope: storageAccount //assigns to storage acct
  name: guid(storageAccount.id, containerGroup.name, applicationInsightsRoleDefinition.id)
  properties: {
    roleDefinitionId: applicationInsightsRoleDefinition.id
    principalId: containerGroup.identity.principalId
    principalType: 'ApplicationInsights'
  }
}
