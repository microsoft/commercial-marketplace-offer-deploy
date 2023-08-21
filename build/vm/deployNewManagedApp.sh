#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <managed_app_version>"
  exit 1
fi

MANAGED_APP_VERSION="$1"

MODM_VM_CLIENT_ID="<CLIENTID>"
MODM_VM_CLIENT_SECRET="<CLIENTSECRET>"
MANAGED_IMAGE_NAME="bobjacnginx$MANAGED_APP_VERSION"
MANAGED_IMAGE_RG_NAME="modmImage"
SUBSCRIPTION_ID="<SUBSCRIPTIONID>"
TENANT_ID="<TENANTID>"

./packer/buildImage.sh $MODM_VM_CLIENT_ID $MODM_VM_CLIENT_SECRET $MANAGED_IMAGE_NAME $MANAGED_IMAGE_RG_NAME $SUBSCRIPTION_ID $TENANT_ID

GALLERY_RESOURCE_GROUP="modmvm"
GALLERY_NAME="bobjacmodm2"
GALLERY_IMAGE_DEFINITION="modm"
GALLERY_IMAGE_VERSION="0.0.$MANAGED_APP_VERSION"
IMAGE="$MANAGED_IMAGE_NAME"
IMAGE_RESOURCE_GROUP="$MANAGED_IMAGE_RG_NAME"
REGION="eastus"

./addImageToGallery.sh $GALLERY_RESOURCE_GROUP $GALLERY_NAME $GALLERY_IMAGE_DEFINITION $GALLERY_IMAGE_VERSION $SUBSCRIPTION_ID $IMAGE $IMAGE_RESOURCE_GROUP $REGION


DEPLOYED_IMAGE_REFERENCE="/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$GALLERY_RESOURCE_GROUP/providers/Microsoft.Compute/galleries/$GALLERY_NAME/images/$GALLERY_IMAGE_DEFINITION/versions/$GALLERY_IMAGE_VERSION"
UIDEF_FILE="createUiDefinition.json"
TEMP_FILE="tempUiDefinition.json"
TEMP_FILE2="tempUiDefinition2.json"

# Use sed to replace <IMAGE_REFERENCE> with the DEPLOYED_IMAGE_REFERENCE
# The -i option has issues with certain platform implementations of the sed command,
# so we use a temporary file for the output and then overwrite the original file
sed "s|<IMAGE_REFERENCE>|$DEPLOYED_IMAGE_REFERENCE|g" "$UIDEF_FILE" > "$TEMP_FILE"

cp -f "$UIDEF_FILE" "./$TEMP_FILE2"
# Replace the original file with the temporary file
cp -f "$TEMP_FILE" "./$UIDEF_FILE"

# Zip up the package for the managed application
./createAppZip.sh

# Update mainTemplate to the origional version with <IMAGE_REFERENCE> so we can rerun next version
cp -f "$TEMP_FILE2" "./$UIDEF_FILE"

# Remove temp files
rm "$TEMP_FILE"
rm "$TEMP_FILE2"

# Create the Service Definition
STORAGE_ACC_RESOURCE_GROUP="bobjacPackageStorage"
STORAGE_ACC_NAME="bobjacmodmscsa"
STORAGE_CONTAINER_NAME="modm$MANAGED_APP_VERSION"
MA_RESOURCE_GROUP="modm$MANAGED_APP_VERSION"
MA_DEFINITION_NAME="modm$MANAGED_APP_VERSION"

./createServiceDefinition.sh $STORAGE_ACC_RESOURCE_GROUP $STORAGE_ACC_NAME $STORAGE_CONTAINER_NAME $MA_RESOURCE_GROUP $MA_DEFINITION_NAME

az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "$DEPLOYED_IMAGE_REFERENCE"