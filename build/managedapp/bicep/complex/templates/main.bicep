param location string
param resourceGroupName string
param artifactsLocationSasToken string
// param resourceGroupName string = 'managedKubeflowRg'
param projectName string = 'managedApp'
var vNetName = '${projectName}BootstrapVnet'
param vNetAddressPrefixes string = '172.18.0.0/16'
var vNetSubnetName = '${projectName}BootstrapSubnet'
param adminUsername string = 'ubuntu'
param vmSize string = 'Standard_D2s_v3'
param vNetSubnetAddressPrefix string = '172.18.0.0/24'
var vmName = '${projectName}BootstrapVm'
var publicIPAddressName = '${projectName}BootstrapVmIp'
var networkInterfaceName = '${projectName}BootstrapVmNic'
// param location string = 'eastus'
var networkSecurityGroupProject = '${projectName}ProjectNsg'
var networkSecurityGroupVnet = '${vNetSubnetName}Nsg'
param clusterName string = 'managedKubeflowCluster'
param dnsPrefix string = 'kubeflow'
param agentCount int = 3
param osDiskSizeGB int = 250

// resource resourceGroup 'Microsoft.Resources/resourceGroups@2022-09-01' = {
//   name: resourceGroupName
//   location: location
//   scope: subscription()
// }

module networkBaseDeployment 'modules/networkBaseDeployment.bicep' = {
  name: 'networkBaseDeployment'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    location: location
    networkSecurityGroupProject: networkSecurityGroupProject
    networkSecurityGroupVnet: networkSecurityGroupVnet
    publicIPAddressName: publicIPAddressName
  }
  dependsOn: [
    // resourceGroup
  ]
}

module virtualNetwork 'modules/virtualNetwork.bicep' = {
  name: 'virtualNetwork'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    location: location
    vNetName: vNetName
    vNetAddressPrefixes: vNetAddressPrefixes
    vNetSubnetName: vNetSubnetName
    vNetSubnetAddressPrefix: vNetSubnetAddressPrefix
    networkSecurityGroupVnet: networkSecurityGroupVnet
  }
  dependsOn: [
    networkBaseDeployment
  ]
}

module networkInterface 'modules/networkInterface.bicep' = {
  name: 'networkInterface'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    location: location
    networkInterfaceName: networkInterfaceName
    publicIPAddressName: publicIPAddressName
    vNetName: vNetName
    vNetSubnetName: vNetSubnetName
  }
  dependsOn: [
    virtualNetwork
  ]
}

module virtualMachine 'modules/virtualMachine.bicep' = {
  name: 'virtualMachine'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    location: location
    vmName: vmName
    adminUsername: adminUsername
    vmSize: vmSize
    networkInterfaceName: networkInterfaceName
    networkInterfaceId: networkInterface.outputs.networkInterfaceId
  }
  dependsOn: [
    networkInterface
  ]
}

module virtualMachinesRoleDefinition 'modules/virtualMachinesRoleDefinition.bicep' = {
  name: 'virtualMachinesRoleDefinition'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    vmName: vmName
    vmPrincipalId: virtualMachine.outputs.vmPrincipalId
  }
  dependsOn: [
    virtualMachine
    AKSCluster
  ]
}

module virtualMachineAKSRoleDefinition 'modules/virtualMachineAKSRoleDefinition.bicep' = {
  name: 'virtualMachineAKSRoleDefinition'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    vmName: vmName
    vmPrincipalId: virtualMachine.outputs.vmPrincipalId
  }
  dependsOn: [
    virtualMachine
    AKSCluster
  ]
}

module virtualMachineExtension 'modules/virtualMachineExtension.bicep' = {
  name: 'virtualMachineExtension'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    location: location
    vmName: vmName
  }
  dependsOn: [
    virtualMachinesRoleDefinition
    virtualMachineAKSRoleDefinition
  ]
}

module AKSCluster 'modules/AKSCluster.bicep' = {
  name: 'AKSCluster'
  scope: az.resourceGroup(resourceGroupName)
  params: {
    location: location
    clusterName: clusterName
    dnsPrefix: dnsPrefix
    nodeResourceGroup: resourceGroupName
    agentCount: agentCount
    vmSize: vmSize
    osDiskSizeGB: osDiskSizeGB
    adminUsername: adminUsername
    vNetName: vNetName
    vNetSubnetName: vNetSubnetName
    vnetSubnetId: virtualNetwork.outputs.vnetSubnetId
  }
  dependsOn: [
    // resourceGroup
    virtualNetwork
  ]
}

output createdResourceGroupName string = resourceGroupName
