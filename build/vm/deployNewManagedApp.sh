#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <managed_app_version>"
  exit 1
fi

# make sure we have a vars file before proceeding
env_pkrvars_file=../../obj/.env.pkrvars

if [ ! -f $env_pkrvars_file ];
then
    echo "./obj/.env.pkrvars file is required."
    exit 1
else
    echo "Packer variables env var file present."
fi

# export packer env variables so they get picked up
export $(grep -v '^#' $env_pkrvars_file | xargs)

MANAGED_APP_VERSION="$1"
export PKR_VAR_managed_image_name=modm-mi-${MANAGED_APP_VERSION}

MODM_VM_CLIENT_ID=$PKR_VAR_client_id
MODM_VM_CLIENT_SECRET=$PKR_VAR_client_secret
MANAGED_IMAGE_NAME=$PKR_VAR_managed_image_name
MANAGED_IMAGE_RG_NAME="modm-dev"
SUBSCRIPTION_ID=$PKR_VAR_subscription_id
TENANT_ID=$PKR_VAR_tenant_id

./packer/buildImage.sh $MODM_VM_CLIENT_ID $MODM_VM_CLIENT_SECRET $MANAGED_IMAGE_NAME $MANAGED_IMAGE_RG_NAME $SUBSCRIPTION_ID $TENANT_ID

GALLERY_RESOURCE_GROUP="modm-dev"
GALLERY_NAME="modm.dev.sig"
GALLERY_IMAGE_DEFINITION="modm"
GALLERY_IMAGE_VERSION="0.0.$MANAGED_APP_VERSION"
IMAGE="$MANAGED_IMAGE_NAME"
IMAGE_RESOURCE_GROUP="$MANAGED_IMAGE_RG_NAME"
REGION="eastus"

./addImageToGallery.sh $GALLERY_RESOURCE_GROUP $GALLERY_NAME $GALLERY_IMAGE_DEFINITION $GALLERY_IMAGE_VERSION $SUBSCRIPTION_ID $IMAGE $IMAGE_RESOURCE_GROUP $REGION

mkdir -p ../../obj

DEPLOYED_IMAGE_REFERENCE="/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$GALLERY_RESOURCE_GROUP/providers/Microsoft.Compute/galleries/$GALLERY_NAME/images/$GALLERY_IMAGE_DEFINITION/versions/$GALLERY_IMAGE_VERSION"
UIDEF_FILE="createUiDefinition.json"
TEMP_FILE="../../obj/createUiDefinition.json"

# Use sed to replace <IMAGE_REFERENCE> with the DEPLOYED_IMAGE_REFERENCE
# The -i option has issues with certain platform implementations of the sed command,
# so we use a temporary file for the output and then overwrite the original file
rm $TEMP_FILE 2> /dev/null
sed "s|<IMAGE_REFERENCE>|$DEPLOYED_IMAGE_REFERENCE|g" "$UIDEF_FILE" > "$TEMP_FILE"

rm ../../obj/mainTemplate.json 2> /dev/null
cp -f mainTemplate.json ../../obj/mainTemplate.json

# Zip up the package for the managed application
./createAppZip.sh


# Create the Service Definition
STORAGE_ACC_RESOURCE_GROUP="modm-dev"
STORAGE_ACC_NAME="modmdev0scsa"
STORAGE_CONTAINER_NAME="modm$MANAGED_APP_VERSION"
MA_RESOURCE_GROUP=$STORAGE_ACC_RESOURCE_GROUP
MA_DEFINITION_NAME="modm$MANAGED_APP_VERSION"

./createServiceDefinition.sh $STORAGE_ACC_RESOURCE_GROUP $STORAGE_ACC_NAME $STORAGE_CONTAINER_NAME $MA_RESOURCE_GROUP $MA_DEFINITION_NAME

az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "$DEPLOYED_IMAGE_REFERENCE"