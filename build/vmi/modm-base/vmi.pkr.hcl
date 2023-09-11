
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

variable "build_resource_group_name" {
  type = string
}

variable "resource_group" {
  type = string
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

source "azure-arm" "base" {
  azure_tags = {
    dept = "Engineering"
    task = "Image deployment"
  }
  build_resource_group_name         = var.build_resource_group_name
  client_id                         = var.client_id
  client_secret                     = var.client_secret
  image_offer                       = "0001-com-ubuntu-server-jammy"
  image_publisher                   = "canonical"
  image_sku                         = "22_04-lts"
  os_type                           = "Linux"
  vm_size                           = "Standard_DS2_v2"
  location                          = var.location
  subscription_id                   = var.subscription_id
  tenant_id                         = var.tenant_id
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
  sources = ["source.azure-arm.base"]

  provisioner "file" {
    sources = [
      "build/container/Dockerfile.jenkins",
      "build/container/Dockerfile.modm",
      "build/container/.dockerignore"
    ]
    destination = "/tmp/"
  }

  provisioner "shell" {
    inline          = [
      # make MODM_HOME env globally available
      "echo MODM_HOME=${var.modm_home} | sudo tee --append /etc/environment",
      "sudo mkdir -p ${var.modm_home}"
    ]
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
    inline = ["/usr/sbin/waagent -force -deprovision+user && export HISTSIZE=0 && sync"]
    inline_shebang  = "/bin/sh -x"
  }

  post-processor "shell-local" {
    inline = [
          "az image delete -n ${var.image_name}-${var.image_version} -g ${var.resource_group} --subscription ${var.subscription_id}"
        ]
  }
}