
@description('Application version in this format: v1.0.0')
param appVersion string

@description('container image')
param containerImage string

@description('Location for all resources.')
param location string = resourceGroup().location

@description('Port to open on the container')
param port int = 8080

@description('The number of CPU cores to allocate to the container.')
param cpuCores int = 1

@description('The amount of memory to allocate to the container in gigabytes.')
param memoryInGb int = 2

@description('The Azure Tenant Id')
param tenantId string

@description('The Azure Subscription Id')
param subscriptionId string

@description('The Azure Resource Group')
param resourceGroupName string

@description('The Azure service bus name')
param serviceBusNamespace string

@description('The email address used for the acme account')
param acmeEmail string

@description('The public http port')
param publicHttpPort int = 80

@description('The public https port')
param publicHttpsPort int = 8443


@description('The behavior of Azure runtime if container has stopped.')
@allowed([
  'Always'
  'Never'
  'OnFailure'
])
param restartPolicy string = 'Always'

var versionSuffix = replace(appVersion, '.', '')

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  name: 'modmstor0${versionSuffix}'
  location: location
  kind: 'StorageV2'
  sku: {
    name: 'Standard_LRS'
  }
}

var fileShareName = '${storageAccount.name}/default/share'
resource fileStore 'Microsoft.Storage/storageAccounts/fileServices/shares@2021-09-01' = {
  name: fileShareName
  properties: {
    shareQuota: 1
    enabledProtocols: 'SMB'
  }
}

var fileShareMountPath = '/opt/share'
var containerName = 'modm-${versionSuffix}'

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2022-10-01-preview' = {
  name: 'modm-group-${versionSuffix}'
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
        name: containerName
        properties: {
          image: containerImage
          ports: [
            {
              port: port
              protocol: 'TCP'
            }
            {
              port: publicHttpPort
              protocol: 'TCP'
            }
            {
              port: publicHttpsPort
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
              mountPath: fileShareMountPath
              readOnly: false
            }
          ]
          environmentVariables: [
            {
              name: 'DB_PATH'
              value: fileShareMountPath
            }
            {
              name: 'AZURE_TENANT_ID'
              value: tenantId
            }
            {
              name: 'AZURE_SUBSCRIPTION_ID'
              value: subscriptionId
            }
            {
              name: 'AZURE_RESOURCE_GROUP'
              value: resourceGroupName
            }
            {
              name: 'AZURE_LOCATION'
              value: location
            }
            {
              name: 'AZURE_SERVICEBUS_NAMESPACE'
              value: serviceBusNamespace
            }
            {
              name: 'ACME_ACCOUNT_EMAIL'
              value: acmeEmail
            }
            {
              name: 'PUBLIC_DOMAIN_NAME'
              value: '${containerName}.${location}.azurecontainer.io'
            }
            {
              name: 'PUBLIC_HTTP_PORT'
              value: string(publicHttpPort)
            }
            {
              name: 'PUBLIC_HTTPS_PORT'
              value: string(publicHttpsPort)
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
          port: publicHttpPort
          protocol: 'TCP'
        }
        {
          port: publicHttpsPort
          protocol: 'TCP'
        }
      ]
      dnsNameLabel: containerName
    }
  }
}

output storageAccountName string = storageAccount.name
output containerGroupName string = containerGroup.name
output containerIPv4Address string = containerGroup.properties.ipAddress.ip
