variable "subscriptionID" {
    type=string
    description = "Azure Subscription ID"
}

variable "tenantID" {
    type=string
    description = "Tenant ID"
}

variable "userAccount" {
    type=string
    description = "User Account Name"
}

variable "computerName" {
    type=string
    description = "Computer Name"
}


variable "vmPasswrd" {
    type=string
    description = "VM Guest OS Password"
}

variable "interface_name" {
    type=list(string)
    description = "interface name"
}

variable "rgName" {
    type=string
    description = "name of RG"
}

variable "location" {
    type=string
    description = "desired region"
}