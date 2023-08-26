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

variable "sig_image_name" {
  type = string
}

variable "sig_image_version" {
  type = string
}

variable "managed_image_name" {
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

source "azure-arm" "base_image" {
  azure_tags = {
    dept = "Engineering"
    task = "Image deployment"
  }
  client_id                         = var.client_id
  client_secret                     = var.client_secret
  image_offer                       = "0001-com-ubuntu-server-jammy"
  image_publisher                   = "canonical"
  image_sku                         = "22_04-lts"
  location                          = "East US"
  managed_image_name                = var.managed_image_name
  managed_image_resource_group_name = "modm-dev"
  os_type                           = "Linux"
  subscription_id                   = var.subscription_id
  tenant_id                         = var.tenant_id
  vm_size                           = "Standard_DS2_v2"
  shared_image_gallery_destination {
      subscription     = var.subscription_id
      resource_group  = var.sig_gallery_resource_group
      gallery_name     = var.sig_gallery_name
      image_name       = var.sig_image_name
      image_version    = var.sig_image_version
      replication_regions = [var.location]
  }
}

build {
  sources = ["source.azure-arm.base_image"]

  provisioner "shell" {
    execute_command = "chmod +x {{ .Path }}; {{ .Vars }} sudo -E sh '{{ .Path }}'"
    inline          = [
      "git clone https://github.com/microsoft/commercial-marketplace-offer-deploy.git /usr/local/source",
      "/usr/local/source/build/vm/scripts/setup.sh",
      "/usr/sbin/waagent -force -deprovision+user && export HISTSIZE=0 && sync",
    ]
    inline_shebang  = "/bin/sh -x"
  }

  post-processor "manifest" {
    output = "manifest.json"
    strip_path = true
    custom_data = {
        source_image_name = "${build.SourceImageName}"
    }
  }
}


