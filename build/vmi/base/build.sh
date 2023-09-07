#!/bin/bash

# import functions
source ./build/vmi/scripts/shared.sh

# -------------------------------------------------------------------------------

# Defaults
RESOURCE_GROUP_NAME="modm-dev-vmi"
IMAGE_NAME="modm-base"
IMAGE_VERSION="0.0.0"
IMAGE_OFFER="modm-base-ubuntu"
BUILD_ENVIRONMENT=$(get_build_environment)


if [ "$BUILD_ENVIRONMENT" = "Local" ]; then
    export_packer_env_vars_from_file
fi

# override default resource group name if packer var declares the resource group
if [ -n "$PKR_VAR_resource_group" ]; then
    RESOURCE_GROUP_NAME=$PKR_VAR_resource_group
fi

source ./build/vmi/scripts/nextversion.sh

if [ $# -ne 1 ]; then
    IMAGE_VERSION=$(get_next_image_version $IMAGE_NAME $RESOURCE_GROUP_NAME)
else
    IMAGE_VERSION="$1"
fi

export PKR_VAR_image_name="${IMAGE_NAME}"
export PKR_VAR_image_version="${IMAGE_VERSION}"

# PREAMBLE
echo ""
echo "Environment Variables "
echo "---------------------------------------------"

echo "BUILD_ENVIRONMENT:        $BUILD_ENVIRONMENT"
echo "RESOURCE_GROUP_NAME:      $RESOURCE_GROUP_NAME"
echo "IMAGE_NAME:               $IMAGE_NAME"
echo "IMAGE_VERSION:            $IMAGE_VERSION"
echo ""
echo "Packer Variables "
echo "---------------------------------------------"
print_packer_variables


# Check if the Shared Image Gallery exists
if ! az sig show --resource-group $RESOURCE_GROUP_NAME --gallery-name $PKR_VAR_image_gallery_name &>/dev/null; then

    # Create a Shared Image Gallery
    az sig create --resource-group $RESOURCE_GROUP_NAME --gallery-name $PKR_VAR_image_gallery_name --location $PKR_VAR_location

    if [ $? -ne 0 ]; then
        echo "Shared Image Gallery creation failed."
        exit 1
    fi
fi

az sig image-definition show \
  --subscription $PKR_VAR_subscription_id \
  --resource-group $RESOURCE_GROUP_NAME \
  --gallery-name $PKR_VAR_image_gallery_name \
  --gallery-image-definition $PKR_VAR_image_name \
  --output none

if [ $? -ne 0 ]; then
  echo "Creating image definition [$PKR_VAR_image_name] in $PKR_VAR_image_gallery_name"
  
  az sig image-definition create \
    --subscription $PKR_VAR_subscription_id \
    --resource-group $RESOURCE_GROUP_NAME \
    --gallery-name $PKR_VAR_image_gallery_name \
    --gallery-image-definition $PKR_VAR_image_name \
    --publisher Microsoft \
    --offer $IMAGE_OFFER \
    --sku $IMAGE_OFFER \
    --os-type Linux \
    --location $PKR_VAR_location

    if [ $? -ne 0 ]; then
        echo "Image definition creation failed."
        exit 1
    fi
fi

echo ""
# Run the Packer command
packer_file=./build/vmi/base/vmi.pkr.hcl
echo "Executing Packer build against [$packer_file]"

# packer init $packer_file
# packer build $packer_file