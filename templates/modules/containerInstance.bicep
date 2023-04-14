@description('Name for the container group')
param name string = 'bobjac26'

@description('Location for all resources.')
param location string = resourceGroup().location

@description('Container image to deploy. Should be of the form repoName/imagename:tag for images stored in public Docker Hub, or a fully qualified URI for other registries. Images from private registries require additional registry credentials.')
param image string = 'bobjac/modm:1.25'

@description('Port to open on the container and the public IP address.')
param port int = 8080

@description('The number of CPU cores to allocate to the container.')
param cpuCores int = 1

@description('The amount of memory to allocate to the container in gigabytes.')
param memoryInGb int = 2

@description('The service principal client ID')
param azureClientId string

@description('The service principal client secret')
@secure()
param azureClientSecret string

@description('The Azure Tenant Id')
param azureTenantId string

@description('The Azure Subscription Id')
param azureSubscriptionId string

@description('The Azure Resource Group')
param azureResourceGroup string

@description('The Azure Location')
param azureLocation string

@description('The Azure Location')
param azureServiceBusNamespace string

@description('The email address used for the acme account')
param acmeEmail string


@description('The behavior of Azure runtime if container has stopped.')
@allowed([
  'Always'
  'Never'
  'OnFailure'
])
param restartPolicy string = 'Always'

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  name: 'inst${name}'
  location: location
  kind: 'StorageV2'
  sku: {
    name: 'Standard_LRS'
  }
}

resource fileStore 'Microsoft.Storage/storageAccounts/fileServices/shares@2021-09-01' = {
  name: '${storageAccount.name}/default/share'
  properties: {
    shareQuota: 1
    enabledProtocols: 'SMB'
  }
}

var containerGroupName = 'installerGroup'
var fqdn = 'dns${name}.${location}.azurecontainer.io'


resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' = {
  name: containerGroupName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    volumes: [
      {
        name: 'filestore'
        azureFile: {
          readOnly: false
          shareName: 'share'
          storageAccountName: storageAccount.name
          storageAccountKey: storageAccount.listKeys().keys[0].value
        }
      }
    ]
    containers: [
      {
        name: 'modm'
        properties: {
          image: image
          ports: [
            {
              port: port
              protocol: 'TCP'
            }
            {
              port: 80
              protocol: 'TCP'
            }
            {
              port: 443
              protocol: 'TCP'
            }
          ]
          resources: {
            requests: {
              cpu: cpuCores
              memoryInGB: memoryInGb
            }
          }
          volumeMounts: [
            {
              name: 'filestore'
              mountPath: '/opt/share'
              readOnly: false
            }
          ]
          environmentVariables: [
            {
              name: 'AZURE_STORAGE_MOUNT_POINT'
              value: '/opt/share'
            }
            {
              name: 'DB_PATH'
              value: '/opt/share'
            }
            { 
              name: 'DB_USE_INMEMORY'
              value: 'false'
            }
            {
              name: 'AZURE_CLIENT_ID'
              value: azureClientId
            }
            {
              name: 'AZURE_TENANT_ID'
              value: azureTenantId
            }
            {
              name: 'AZURE_CLIENT_SECRET'
              value: azureClientSecret
            }
            {
              name: 'AZURE_SUBSCRIPTION_ID'
              value: azureSubscriptionId
            }
            {
              name: 'AZURE_RESOURCE_GROUP'
              value: azureResourceGroup
            }
            {
              name: 'AZURE_LOCATION'
              value: azureLocation
            }
            {
              name: 'AZURE_SERVICEBUS_NAMESPACE'
              value: azureServiceBusNamespace
            }
            {
              name: 'ACME_ACCOUNT_EMAIL'
              value: acmeEmail
            }
            {
              name: 'INSTALLER_DOMAIN_NAME'
              value: fqdn
            }
          ]
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: restartPolicy
    ipAddress: {
      type: 'Public'
      ports: [
        {
          port: port
          protocol: 'TCP'
        }
        {
          port: 80
          protocol: 'TCP'
        }
        {
          port: 443
          protocol: 'TCP'
        }
      ]
      dnsNameLabel: 'dns${name}'
    }
  }
}

output containerIPv4Address string = containerGroup.properties.ipAddress.ip
