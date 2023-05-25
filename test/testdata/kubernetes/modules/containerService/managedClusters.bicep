param name string = ''
param serviceCidr string
param dnsServiceIP string
param podCidr string

param nodePool_count int = 1
@allowed([
  'System'
  'User'
])
param nodePool_mode string = 'System'
param nodePool_vmSize string = 'Standard_D4s_v3'
param enablePrivateCluster bool = true
param identity object = {
  type: 'SystemAssigned'
}
param podIdentityProfile object = {}

param location string = ''
param logAnalyticsWorkspaceResourceID string = ''
param nodeResourceGroup string = ''
param privateDnsZoneId string = ''
param vnetSubnetID string = ''

var dnsPrefix = name

var defaultProperties = {
  kubernetesVersion: '1.25.5'
  nodeResourceGroup: nodeResourceGroup
  enableRBAC: true
  dnsPrefix: dnsPrefix
  networkProfile: {
    serviceCidr: serviceCidr
    dnsServiceIP: dnsServiceIP
    podCidr: podCidr
    dockerBridgeCidr: '172.17.0.0/24'
  }
  agentPoolProfiles: [
    {
      name: 'agentpool'
      count: nodePool_count
      type: 'VirtualMachineScaleSets'
      enableAutoScaling: true
      vmSize: nodePool_vmSize
      osType: 'Linux'
      maxCount: 20
      minCount: 1
      mode: nodePool_mode
      vnetSubnetID: vnetSubnetID
      kubeletConfig: {
        containerLogMaxFiles: 5
        containerLogMaxSizeMB: 500
      }
    }
  ]
  apiServerAccessProfile: {
    enablePrivateCluster: enablePrivateCluster
    privateDNSZone: privateDnsZoneId
  }
  podIdentityProfile: podIdentityProfile
}
var addonProfiles = !empty(logAnalyticsWorkspaceResourceID) ? {
  addonProfiles: {
    omsagent: {
      enabled: true
      config: {
        logAnalyticsWorkspaceResourceID: logAnalyticsWorkspaceResourceID
      }
    }
  }
} : {}
var properties = union(defaultProperties, addonProfiles)

resource managedCluster 'Microsoft.ContainerService/managedClusters@2022-07-01' = {
  name: name
  location: location
  identity: identity
  properties: properties
}

output name string = managedCluster.name
