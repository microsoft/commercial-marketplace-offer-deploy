
var location = 'eastus2'

module kubernetes 'modules/kubernetes.bicep' = {
  name: 'kubernetes'
  params: {
    dnsPrefix: 'bobjac'
    location: location
    linuxAdminUsername: 'bobjac'
    sshRSAPublicKey: '[INSERT SSH KEY HERE]'
  }
}
