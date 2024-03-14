param vmName string
param vmPrincipalId string

resource roleAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  name: guid(resourceGroup().id, 'Owner', vmName)
  properties: {
    roleDefinitionId: concat(subscription().id, '/providers/Microsoft.Authorization/roleDefinitions/', '8e3af657-a8ff-443c-a75c-2fe8c4bcb635')
    principalId: vmPrincipalId
    scope: resourceGroup().id
  }
}
