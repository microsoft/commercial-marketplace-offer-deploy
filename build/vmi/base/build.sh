#!/bin/bash

BASE_VERSION="0.0.0"
source ./build/vmi/scripts/nextversion.sh

if [ $# -ne 1 ]; then
    BASE_VERSION=$(get_next_image_version "modmvmi-base" "modm-dev-vmi")
else
    BASE_VERSION="$1"
fi

echo "Building modm version $BASE_VERSION"

# Check if running in GitHub Actions environment
if [ -n "$GITHUB_ACTIONS" ]; then
    echo "Running in GitHub Actions environment"
    # You don't need to check for the .env.pkrvars file or export variables here
else
    echo "Running locally"
    # make sure we have a vars file before proceeding
    env_pkrvars_file=./obj/.env.pkrvars

    if [ ! -f $env_pkrvars_file ];
    then
        echo "./obj/.env.pkrvars file is required."
        exit 1
    else
        echo "Packer variables env var file present."
    fi

    # export packer env variables so they get picked up
    export $(grep -v '^#' $env_pkrvars_file | xargs)
fi

export PKR_VAR_sig_image_version=modmvmi-base-${BASE_VERSION}
export PKR_VAR_managed_image_name=modmvmi-base-${BASE_VERSION}
export PKR_VAR_sig_image_version=${BASE_VERSION}

# Check if the Shared Image Gallery exists
if ! az sig show --resource-group $PKR_VAR_sig_gallery_resource_group --gallery-name $PKR_VAR_sig_gallery_name &>/dev/null; then

    # Create a Shared Image Gallery
    az sig create --resource-group $PKR_VAR_sig_gallery_resource_group --gallery-name $PKR_VAR_sig_gallery_name --location $PKR_VAR_location

    echo "Shared Image Gallery created."
fi

az sig image-definition show \
  --subscription $PKR_VAR_subscription_id \
  --resource-group $PKR_VAR_sig_gallery_resource_group \
  --gallery-name $PKR_VAR_sig_gallery_name \
  --gallery-image-definition $PKR_VAR_sig_image_name \
  --output none

if [ $? -eq 0 ]; then
  echo "Image definition $image_definition_name already exists."
else
  echo "Image definition $image_definition_name doesn't exist. Creating..."
  
  az sig image-definition create \
    --subscription $PKR_VAR_subscription_id \
    --resource-group $PKR_VAR_sig_gallery_resource_group \
    --gallery-name $PKR_VAR_sig_gallery_name \
    --gallery-image-definition $PKR_VAR_sig_image_name \
    --publisher Microsoft \
    --offer MODMBaseVMI \
    --sku MODMBaseVMI \
    --os-type Linux
fi


# Run the Packer command
packer init ./build/vmi/base/base.pkr.hcl
packer build ./build/vmi/base/base.pkr.hcl