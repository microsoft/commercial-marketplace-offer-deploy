param name string = ''
param location string = ''

@allowed([
  'Premium_LRS'
  'Premium_ZRS'
  'Standard_GRS'
  'Standard_GZRS'
  'Standard_LRS'
  'Standard_RAGRS'
  'Standard_RAGZRS'
  'Standard_ZRS'
])
param sku_name string = 'Standard_LRS'

@allowed([
  'BlobStorage'
  'BlockBlobStorage'
  'FileStorage'
  'Storage'
  'StorageV2'
])
param kind string = 'StorageV2'

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-08-01' = {
  name: name
  location: location
  sku: {
    name: sku_name
  }
  kind: kind
  properties: {
    allowBlobPublicAccess: false
    publicNetworkAccess: 'Disabled'
    minimumTlsVersion: 'TLS1_2'
    networkAcls: {
      defaultAction: 'Deny'
    }
  }
}
resource storageAccount_blob 'Microsoft.Storage/storageAccounts/blobServices@2021-08-01' = {
  parent: storageAccount
  name: 'default'
  properties: {
    deleteRetentionPolicy: {
      enabled: true
      days: 181
    }
    isVersioningEnabled: true
    changeFeed: {
      enabled: true
    }
    restorePolicy: {
      enabled: true
      days: 180
    }
    containerDeleteRetentionPolicy: {
      enabled: true
      days: 181
    }
  }
}
output name string = storageAccount.name
output id string = storageAccount.id
