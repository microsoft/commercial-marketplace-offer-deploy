
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

variable "resource_group" {
  type = string
}

variable "base_image_name" {
  type = string
}

variable "base_image_version" {
  type = string
  default = "latest"
}

variable "image_name" {
  type = string
}

variable "image_version" {
  type = string
}

variable "image_gallery_name" {
  type = string
}

variable "modm_home" {
  type = string
  default = "/usr/local/modm"
}

variable "modm_repo_branch" {
  type = string
  default = "main"
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
  location                          = var.location
  subscription_id                   = var.subscription_id
  tenant_id                         = var.tenant_id
  os_type                           = "Linux"
  vm_size                           = "Standard_DS2_v2"

  # defines the base image source
  shared_image_gallery {
    subscription     = var.subscription_id
    resource_group   = var.resource_group
    gallery_name     = var.image_gallery_name
    image_name       = var.base_image_name
    image_version    = var.base_image_version
  }

  managed_image_name                = "${var.image_name}-${var.image_version}"
  managed_image_resource_group_name = var.resource_group

  shared_image_gallery_destination {
    subscription     = var.subscription_id
    resource_group   = var.resource_group
    gallery_name     = var.image_gallery_name
    image_name       = var.image_name
    image_version    = var.image_version
    replication_regions = [var.location]
  }
}

build {
  sources = ["source.azure-arm.modm"]

  provisioner "file" {
    sources = [
      "build/vmi/${var.image_name}/files/modm.service",
      "build/vmi/${var.image_name}/files/Caddyfile",
      "build/vmi/${var.image_name}/files/docker-compose.yml",
    ]
    destination = "/tmp/"
  }

  provisioner "shell" {
    environment_vars = [
        "MODM_HOME=${var.modm_home}",
        "MODM_REPO_BRANCH=${var.modm_repo_branch}"
      ]
    script = "build/vmi/${var.image_name}/setup.sh"
  }

  provisioner "shell" {
    execute_command = "chmod +x {{ .Path }}; {{ .Vars }} sudo -E sh '{{ .Path }}'"
    inline          = [
      "/usr/sbin/waagent -force -deprovision+user && export HISTSIZE=0 && sync",
    ]
    inline_shebang  = "/bin/sh -x"
  }

  post-processor "shell-local" {
    inline = [
          "az image delete -n ${var.image_name}-${var.image_version} -g ${var.resource_group} --subscription ${var.subscription_id}"
        ]
  }
}
