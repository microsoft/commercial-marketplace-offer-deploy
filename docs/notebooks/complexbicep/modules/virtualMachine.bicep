param location string
param vmName string
param adminUsername string
param vmSize string
param networkInterfaceName string
param networkInterfaceId string

resource virtualMachine 'Microsoft.Compute/virtualMachines@2021-11-01' = {
  name: vmName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    hardwareProfile: {
      vmSize: vmSize
    }
    osProfile: {
      computerName: vmName
      adminUsername: adminUsername
      adminPassword: 'EJYtrSN34qpVbS9jKCWdKFoz1NmeA58RYvCbQfGn0FHHBlVTJQ46Gkuh7edB2Lfo'
    }
    storageProfile: {
      imageReference: {
        publisher: 'Canonical'
        offer: '0001-com-ubuntu-server-jammy'
        sku: '22_04-lts-gen2'
        version: 'latest'
      }
      osDisk: {
        createOption: 'fromImage'
      }
    }
    networkProfile: {
      networkInterfaces: [
        {
          id: networkInterfaceId
        }
      ]
    }
  }
}

output vmPrincipalId string = virtualMachine.identity.principalId
