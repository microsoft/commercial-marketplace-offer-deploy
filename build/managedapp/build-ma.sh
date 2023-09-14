#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <managed_app_version>"
  exit 1
fi

MANAGED_APP_VERSION="$1"
echo "The current directory is: $(pwd)"

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
cp -f ./build/managedapp/content.zip ./obj/content.zip

echo "The ./obj directory contains: $(ls -la ./obj)"

# Zip up the package for the managed application
./build/managedapp/createAppZip.sh


# Create the Service Definition
STORAGE_ACC_RESOURCE_GROUP=$MANAGED_APP_STORAGE_RG
STORAGE_ACC_NAME=$MANAGED_APP_STORAGE_NAME
BUILDNUM=$(echo "$MANAGED_APP_VERSION" | awk -F. '{print $3}')
STORAGE_CONTAINER_NAME="modm$BUILDNUM"
MA_RESOURCE_GROUP=$STORAGE_ACC_RESOURCE_GROUP
MA_DEFINITION_NAME="modm$MANAGED_APP_VERSION"

./build/managedapp/createServiceDefinition.sh $STORAGE_ACC_RESOURCE_GROUP $STORAGE_ACC_NAME $STORAGE_CONTAINER_NAME $MA_RESOURCE_GROUP $MA_DEFINITION_NAME

az role assignment create --assignee 1cf33839-e2dd-49a4-a41f-03a52b70a203 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "$DEPLOYED_IMAGE_REFERENCE"