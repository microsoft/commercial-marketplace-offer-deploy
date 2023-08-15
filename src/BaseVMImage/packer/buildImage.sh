#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 6 ]; then
  echo "Usage: $0 <client_id> <client_secret> <managed_image_name> <managed_image_resource_group_name> <subscription_id> <tenant_id>"
  exit 1
fi

# Assign parameters to variables
CLIENT_ID="$1"
CLIENT_SECRET="$2"
MANAGED_IMAGE_NAME="$3"
MANAGED_IMAGE_RG_NAME="$4"
SUBSCRIPTION_ID="$5"
TENANT_ID="$6"

# Update Packer template file with parameter values
cat <<EOF > modm.pkr.hcl
source "azure-arm" "autogenerated_1" {
  azure_tags = {
    dept = "Engineering"
    task = "Image deployment"
  }
  client_id                         = "$CLIENT_ID"
  client_secret                     = "$CLIENT_SECRET"
  image_offer                       = "0001-com-ubuntu-server-jammy"
  image_publisher                   = "canonical"
  image_sku                         = "22_04-lts"
  location                          = "East US"
  managed_image_name                = "$MANAGED_IMAGE_NAME"
  managed_image_resource_group_name = "$MANAGED_IMAGE_RG_NAME"
  os_type                           = "Linux"
  subscription_id                   = "$SUBSCRIPTION_ID"
  tenant_id                         = "$TENANT_ID"
  vm_size                           = "Standard_DS2_v2"
}

build {
  sources = ["source.azure-arm.autogenerated_1"]

  provisioner "shell" {
    execute_command = "chmod +x {{ .Path }}; {{ .Vars }} sudo -E sh '{{ .Path }}'"
    inline          = ["apt-get update", "apt-get upgrade -y", "apt-get -y install nginx", "git clone https://github.com/microsoft/commercial-marketplace-offer-deploy.git /usr/local/source", "/usr/sbin/waagent -force -deprovision+user && export HISTSIZE=0 && sync"]
    inline_shebang  = "/bin/sh -x"
  }
}
EOF

# Run the Packer command
packer build modm.pkr.hcl

# Create role assignment
az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$MANAGED_IMAGE_RG_NAME/providers/Microsoft.Compute/images/$MANAGED_IMAGE_NAME"