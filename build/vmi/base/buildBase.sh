#!/bin/bash

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

BASE_VERSION="$1"
export PKR_VAR_sig_image_version=modm-base-${BASE_VERSION}
export PKR_VAR_managed_image_name=modm-base-${BASE_VERSION}
export PKR_VAR_sig_image_version=${BASE_VERSION}

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
    --offer MODMBase \
    --sku MODMBase \
    --os-type Linux

fi


# Run the Packer command
packer init ./build/vmi/base/base.pkr.hcl
packer build ./build/vmi/base/base.pkr.hcl