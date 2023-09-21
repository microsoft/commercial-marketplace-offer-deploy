terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "=2.99.0"
    }
  }
}


provider "azurerm" {
  subscription_id = var.subscriptionID
  tenant_id       = var.tenantID
  features {
    }
}

module "build_eastRG" {
  source = "./modules/ResourceGroup"
  rg_Name = var.resourceGrp1
  rg_Location = var.resourceGrp1_Location
  tenant_ID = var.tenantID
  subscription_ID = var.subscriptionID
}

module "build_westRG" {
  source = "./modules/ResourceGroup"
  rg_Name = var.resourceGrp2
  rg_Location = var.resourceGrp2_Location
  tenant_ID = var.tenantID
  subscription_ID = var.subscriptionID
}

module "build_EastVNet" {
    source = "./modules/VNet"
    rgName = module.build_eastRG.rg_name_out
    location = module.build_eastRG.rg_location_out
    subscriptionID = var.subscriptionID
    tenantID = var.tenantID
    Vnet_Name = var.vnet1_Name
    Vnet_AddressSpace = ["10.0.0.0/16"]
    SubnetPrefix = ["10.0.0.0/24"]
    NSG_Name = var.nsg1_Name
    public_IP = var.pub1_Name
    interface_name = var.nic1_name
}

module "build_WestVNet" {
    source = "./modules/VNet"
    rgName = module.build_westRG.rg_name_out
    location = module.build_westRG.rg_location_out
    subscriptionID = var.subscriptionID
    tenantID = var.tenantID
    Vnet_Name = var.vnet2_Name
    Vnet_AddressSpace = ["10.5.0.0/16"]
    SubnetPrefix = ["10.5.0.0/24"]
    NSG_Name = var.nsg2_Name
    public_IP = var.pub2_Name
    interface_name = var.nic2_name
}

module "build_peerWest" {
    source = "./modules/Peering"
    subscriptionID = var.subscriptionID
    tenantID = var.tenantID
    rgName = module.build_westRG.rg_name_out
    peerName = var.west_PeerName
    localPeerVnetName = module.build_WestVNet.networkName_out
    remotePeerVnetID = module.build_EastVNet.networkVNetID_out
}

module "build_peerEast" {
    source = "./modules/Peering"
    subscriptionID = var.subscriptionID
    tenantID = var.tenantID
    rgName = module.build_eastRG.rg_name_out
    peerName = var.east_PeerName
    localPeerVnetName = module.build_EastVNet.networkName_out
    remotePeerVnetID = module.build_WestVNet.networkVNetID_out
}


module "build_EastVM" {
    source = "./modules/VM"
    rgName = module.build_eastRG.rg_name_out
    location = module.build_eastRG.rg_location_out
    subscriptionID = var.subscriptionID
    tenantID = var.tenantID
    userAccount = var.vm_User
    computerName = var.vmGuestOS1_Name
    vmPasswrd = var.vm_Password
    interface_name = [module.build_EastVNet.networkID_out]      
}

module "build_westVM" {
    source = "./modules/VM"
    rgName = module.build_westRG.rg_name_out
    location = module.build_westRG.rg_location_out
    subscriptionID = var.subscriptionID
    tenantID = var.tenantID
    userAccount = var.vm_User
    computerName = var.vmGuestOS2_Name
    vmPasswrd = var.vm_Password
    interface_name = [module.build_WestVNet.networkID_out]
}


