@description('Name for the container group')
param name string = 'bobjacharness48'

@description('Location for all resources.')
param location string = resourceGroup().location

@description('Container image to deploy. Should be of the form repoName/imagename:tag for images stored in public Docker Hub, or a fully qualified URI for other registries. Images from private registries require additional registry credentials.')
param image string = 'bobjac/modmtestharness:1.35'

@description('Port to open on the container and the public IP address.')
param port int = 8280

@description('The number of CPU cores to allocate to the container.')
param cpuCores int = 1

@description('The amount of memory to allocate to the container in gigabytes.')
param memoryInGb int = 2

@description('The api endpoint')
param apiEndpoint string

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


@description('The behavior of Azure runtime if container has stopped.')
@allowed([
  'Always'
  'Never'
  'OnFailure'
])
param restartPolicy string = 'Always'

var containerGroupName = 'harnessGroup'
var fqdn = 'dns${name}.${location}.azurecontainer.io'

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' = {
  name: containerGroupName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    containers: [
      {
        name: 'modmtestharness'
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
          environmentVariables: [
            {
              name: 'MODM_API_ENDPOINT'
              value: apiEndpoint
            }
            {
              name: 'API_URI'
              value: apiEndpoint
            }
            {
              name: 'MODM_SUBSCRIPTION'
              value: azureSubscriptionId
            }
            {
              name: 'MODM_RESOURCE_GROUP'
              value: azureResourceGroup
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
              name: 'TEMPLATE_PATH'
              value: '/mainTemplateBicep.json'
            }
            {
              name: 'TEMPLATEPARAMS_PATH'
              value: '/parametersBicep.json'
            }
            {
              name: 'CALLBACK'
              value: 'dns${name}'
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
