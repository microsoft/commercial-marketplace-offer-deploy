variable "client_id" {
  type = string
}

variable "client_secret" {
  type = string
}

variable "subscription_id" {
  type = string
}

variable "tenant_id" {
  type = string
}

variable "location" {
  type = string
}

variable "managed_image_name" {
  type = string
}

variable "resource_group_name" {
  type = string
}

packer {
  required_plugins {
    azure = {
      version = ">= 2.0.0"
      source  = "github.com/hashicorp/azure"
    }
  }
}

source "azure-arm" "modm" {
  azure_tags = {
    dept = "Engineering"
    task = "Image deployment"
  }
  client_id                         = var.client_id
  client_secret                     = var.client_secret
  image_offer                       = "0001-com-ubuntu-server-jammy"
  image_publisher                   = "canonical"
  image_sku                         = "22_04-lts"
  location                          = var.location
  managed_image_name                = var.managed_image_name
  managed_image_resource_group_name = var.resource_group_name
  os_type                           = "Linux"
  subscription_id                   = var.subscription_id
  tenant_id                         = var.tenant_id
  vm_size                           = "Standard_DS2_v2"
}

build {
  sources = ["source.azure-arm.modm"]

  provisioner "shell" {
    execute_command = "chmod +x {{ .Path }}; {{ .Vars }} sudo -E sh '{{ .Path }}'"
    inline          = [
      "git clone https://github.com/microsoft/commercial-marketplace-offer-deploy.git /usr/local/source",
      "/usr/local/source/build/vmi/scripts/setup.sh",
      "/usr/sbin/waagent -force -deprovision+user && export HISTSIZE=0 && sync",
    ]
    inline_shebang  = "/bin/sh -x"
  }
}
