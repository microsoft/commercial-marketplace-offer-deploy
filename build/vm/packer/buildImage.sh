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

# export packer env variables so they get picked up
export $(grep -v '^#' ../../bin/.env.packer | xargs)

# Run the Packer command
packer init ./modm.pkr.hcl
packer build ./modm.pkr.hcl

# Create role assignment
az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$MANAGED_IMAGE_RG_NAME/providers/Microsoft.Compute/images/$MANAGED_IMAGE_NAME"
