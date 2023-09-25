variable "subscriptionID" {
    type=string
    description = "Azure Subscription ID"
}

variable "tenantID" {
    type=string
    description = "Tenant ID"
}

variable "rgName" {
    type=string
    description = "name of RG"
}

variable "localPeerVnetName" {
    type=string
    description = "Source vNet name for Peer"
}

variable "remotePeerVnetID" {
    type=string
    description = "remote VNet ID for Peer"
}

variable "peerName" {
    type=string
    description = "PeerName"
}
