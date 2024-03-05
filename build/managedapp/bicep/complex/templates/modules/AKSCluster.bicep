param location string
param clusterName string
param dnsPrefix string
param nodeResourceGroup string
param agentCount int
param vmSize string
param osDiskSizeGB int
param adminUsername string
param vNetName string
param vNetSubnetName string
param vnetSubnetId string

resource aksCluster 'Microsoft.ContainerService/managedClusters@2021-07-01' = {
  name: clusterName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    dnsPrefix: dnsPrefix
    nodeResourceGroup: nodeResourceGroup
    agentPoolProfiles: [
      {
        name: 'agentpool'
        osDiskSizeGB: osDiskSizeGB
        count: agentCount
        vmSize: vmSize
        osType: 'Linux'
        mode: 'System'
        vnetSubnetID: vnetSubnetId
      }
    ]
    linuxProfile: {
      adminUsername: adminUsername
      ssh: {
        publicKeys: [
          {
            keyData: 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDwtHS+ABlV5S8jxubzS6IZZY2dywMad5qL68aKplglt4YyFFAYAbTO9L+ZLvrJCzlBSnhkWOrEeGNPjokbpLVxiIv9/Ma/CHDVdG2rdPsns7b6yH0Bjnp7AAb5IJsY+3t2OHhNLGy08iTPqk8rJaQa4F2w/PCO7q4Qz+e8xzxEfvmF5oFqUwjrXuaS8JdDUEyjtwidjn0eW4Z/Qeg4pqKbHu/w17grX7SaZ8pfoGBT9Kc5Z4zia4djzLHVcULc32pEuONHl9yIg2waR0+rFQs4YhtneJDJPwGCH2FmXG11IaINUHT0d3KKxVx1U9YmoDT1rZxacloA0icwKl08uHOB'
          }
        ]
      }
    }
  }
}
