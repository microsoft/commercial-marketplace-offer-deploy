#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <managed_app_version>"
  exit 1
fi

MANAGED_APP_VERSION="$1"

mkdir -p ./obj

DEPLOYED_IMAGE_REFERENCE="/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$GALLERY_RESOURCE_GROUP/providers/Microsoft.Compute/galleries/$GALLERY_NAME/images/$GALLERY_IMAGE_DEFINITION/versions/$GALLERY_IMAGE_VERSION"
UIDEF_FILE="./build/managedapp/createUiDefinition.json"
TEMP_FILE="./obj/createUiDefinition.json"

# Use sed to replace <IMAGE_REFERENCE> with the DEPLOYED_IMAGE_REFERENCE
# The -i option has issues with certain platform implementations of the sed command,
# so we use a temporary file for the output and then overwrite the original file
rm $TEMP_FILE 2> /dev/null
sed "s|<IMAGE_REFERENCE>|$DEPLOYED_IMAGE_REFERENCE|g" "$UIDEF_FILE" > "$TEMP_FILE"

rm ./obj/mainTemplate.json 2> /dev/null
cp -f ./build/managedapp/mainTemplate.json ./obj/mainTemplate.json

# Zip up the package for the managed application
./build/managedapp/createAppZip.sh


# Create the Service Definition
STORAGE_ACC_RESOURCE_GROUP=$MANAGED_APP_STORAGE_RG
STORAGE_ACC_NAME=$MANAGED_APP_STORAGE_NAME
STORAGE_CONTAINER_NAME="modm$MANAGED_APP_VERSION"
MA_RESOURCE_GROUP=$STORAGE_ACC_RESOURCE_GROUP
MA_DEFINITION_NAME="modm$MANAGED_APP_VERSION"

./build/managedapp/createServiceDefinition.sh $STORAGE_ACC_RESOURCE_GROUP $STORAGE_ACC_NAME $STORAGE_CONTAINER_NAME $MA_RESOURCE_GROUP $MA_DEFINITION_NAME

az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "$DEPLOYED_IMAGE_REFERENCE"