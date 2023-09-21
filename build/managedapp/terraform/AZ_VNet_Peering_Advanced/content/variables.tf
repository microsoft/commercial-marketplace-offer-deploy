#Required Input -- Enter your Azure Subscription ID#
variable "subscriptionID" {
    type=string
    default = ""
}


#Required Input -- Enter your Azure AD Tenant ID#
variable "tenantID" {
    type=string
    default = ""
}

#Required Input -- Enter your password for Guest OS access#
#Important Note: this is not secure and recommended to use Azure Key Vault#
variable "vm_Password" {
    type=string
    default = "GPSCodeWith123"
}

variable "resourceGrp1" {
    type=string
    default = "ContosoEast"
}

variable "resourceGrp2" {
    type=string
    default = "ContosoWest"
}

variable "resourceGrp1_Location" {
    type=string
    default = "East US"
}

variable "resourceGrp2_Location" {
    type=string
    default = "West US"
}

variable "vnet1_Name" {
    type=string
    default = "eastNet"
}

variable "vnet2_Name" {
    type=string
    default = "westNet"
}


variable "nsg1_Name" {
    type=string
    default = "eastNSG"
}

variable "nsg2_Name" {
    type=string
    default = "westNSG"
}

variable "nic1_name" {
    type=string
    default = "eastNic" 
}

variable "nic2_name" {
    type=string
    default = "westNic" 
}

variable "pub1_Name" {
    type=string
    default = "eastPubIP"
}

variable "pub2_Name" {
    type=string
    default = "westPubIP"
}


variable "vm_User" {
    type=string
    default = "demousr"
}



variable "vmGuestOS1_Name" {
    type=string
    default = "contosoVM1"
}

variable "vmGuestOS2_Name" {
    type=string
    default = "contosoVM2"
}

variable "east_PeerName" {
    type=string
    default = "peerEast2West"
}

variable "west_PeerName" {
    type=string
    default = "peerWest2East"
}