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

variable "sig_gallery_resource_group" {
  type = string
}

variable "sig_gallery_name" {
  type = string
}

variable "sig_image_name_modm" {
  type = string
}

variable "sig_image_version_modm" {
  type = string
}

variable "managed_image_name" {
  type = string
}

variable "managed_image_name_modm" {
  type = string
}

variable "managed_image_resource_group_modm" {
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

source "azure-arm" "modm_image" {
  azure_tags = {
    dept = "Engineering"
    task = "Image deployment"
  }
  client_id                         = var.client_id
  client_secret                     = var.client_secret
  location                          = var.location
  managed_image_name                = var.managed_image_name_modm
  managed_image_resource_group_name = var.managed_image_resource_group_modm
  os_type                           = "Linux"
  subscription_id                   = var.subscription_id
  tenant_id                         = var.tenant_id
  vm_size                           = "Standard_DS2_v2"
  custom_managed_image_name         = "modm-base-0.0.3"
  custom_managed_image_resource_group_name = var.sig_gallery_resource_group
  shared_image_gallery_destination {
      subscription     = var.subscription_id
      resource_group  = var.sig_gallery_resource_group
      gallery_name     = var.sig_gallery_name
      image_name       = var.sig_image_name_modm
      image_version    = var.sig_image_version_modm
      replication_regions = [var.location]
  }
}

build {
  sources = ["source.azure-arm.modm_image"]

  provisioner "shell" {
    execute_command = "chmod +x {{ .Path }}; {{ .Vars }} sudo -E sh '{{ .Path }}'"
    inline          = [
      "git clone https://github.com/microsoft/commercial-marketplace-offer-deploy.git /usr/local/modmsource",
      "/usr/local/modmsource/build/vm/scripts/build.sh",
      "/usr/sbin/waagent -force -deprovision+user && export HISTSIZE=0 && sync",
    ]
    inline_shebang  = "/bin/sh -x"
  }
}
