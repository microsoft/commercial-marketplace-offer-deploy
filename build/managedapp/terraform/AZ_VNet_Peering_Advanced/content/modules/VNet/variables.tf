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

variable "location" {
    type=string
    description = "Region"
}

variable "Vnet_Name" {
    type=string
    description = "Name of Vnet"
}

variable "Vnet_AddressSpace" {
    type = list(string)
    description = "Address Space for VNet"
}

variable "SubnetPrefix" {
    type= list(string)
    description = "Subnet Prefix for default subnet"
}

variable "NSG_Name" {
    type=string
    description = "Name of Network Security Group"
}

variable "public_IP" {
    type=string
    description = "Name of Public IP"
}

variable "interface_name" {
    type=string
    description = "interface name"
}




