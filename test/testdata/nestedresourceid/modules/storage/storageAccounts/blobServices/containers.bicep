param name string = ''
param parent string = ''

var containerName = '${parent}/default/${name}'
resource container 'Microsoft.Storage/storageAccounts/blobServices/containers@2022-05-01' = {
  name: containerName
}
