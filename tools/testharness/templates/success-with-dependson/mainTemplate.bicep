

param suffix string = substring(uniqueString(utcNow()), 0, 5)

var storageAccountName = toLower(format('{0}{1}', 'storageAccounts-', suffix))

resource storageAccount 'Microsoft.Storage/storageAccounts@2022-09-01' = {
  name: storageAccountName
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
   properties: {
     accessTier: 'Hot'
   }
  tags: resourceTags
}
