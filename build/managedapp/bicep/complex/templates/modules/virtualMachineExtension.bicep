param location string
param vmName string

resource customScriptExtension 'Microsoft.Compute/virtualMachines/extensions@2020-06-01' = {
  name: '${vmName}/customScript'
  location: location
  properties: {
    publisher: 'Microsoft.Azure.Extensions'
    type: 'CustomScript'
    typeHandlerVersion: '2.1'
    settings: {
      script: 'IyEvYmluL2Jhc2gKIyBpbnN0YWxsIHRhaWxzY2FsZSBhbmQgam9pbiBuZXR3b3JrCmN1cmwgLWZzU0wgaHR0cHM6Ly90YWlsc2NhbGUuY29tL2luc3RhbGwuc2ggfCBzaAo='
    }
  }
}
