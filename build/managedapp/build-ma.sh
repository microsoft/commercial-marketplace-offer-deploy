#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 3 ]; then
  echo "Usage: $0 <managed_app_version> <deployed_image_reference> <scenario>"
  exit 1
fi

MANAGED_APP_VERSION="$1"

# in the format of {deploymentType}/{name}
SCENARIO="$3"
SCENARIO=${SCENARIO//$'\n'/}
SCENARIO=${SCENARIO//$'\r'/}

echo "The scenario: $SCENARIO"

echo "creating directories in $(pwd)"
mkdir -p ./obj
mkdir -p ./bin

#DEPLOYED_IMAGE_REFERENCE="/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$GALLERY_RESOURCE_GROUP/providers/Microsoft.Compute/galleries/$GALLERY_NAME/images/$GALLERY_IMAGE_DEFINITION/versions/$GALLERY_IMAGE_VERSION"
DEPLOYED_IMAGE_REFERENCE="$2"
echo "The deployed image reference is: $DEPLOYED_IMAGE_REFERENCE"
UIDEF_FILE="./build/managedapp/$SCENARIO/createUiDefinition.json"
TEMP_FILE="./obj/createUiDefinition.json"

echo "The UIDEF_FILE is: $UIDEF_FILE"

# Assign the Reader role to the Managed Application Service Principal
az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "$DEPLOYED_IMAGE_REFERENCE"

# Use sed to replace <IMAGE_REFERENCE> with the DEPLOYED_IMAGE_REFERENCE
# The -i option has issues with certain platform implementations of the sed command,
# so we use a temporary file for the output and then overwrite the original file
rm $TEMP_FILE 2> /dev/null

echo "prior to sed command, the UIDEF_FILE contains: $(cat $UIDEF_FILE)"

sed "s|<IMAGE_REFERENCE>|$DEPLOYED_IMAGE_REFERENCE|g" "$UIDEF_FILE" > "$TEMP_FILE"

rm ./obj/mainTemplate.json 2> /dev/null
rm ./obj/viewDefinition.json 2> /dev/null
rm ./obj/viewDefinition.json 2> /dev/null

cp -f ./build/managedapp/$SCENARIO/mainTemplate.json ./obj/mainTemplate.json
cp -f ./build/managedapp/$SCENARIO/viewDefinition.json ./obj/viewDefinition.json

echo "The ./obj directory contains: $(ls -la ./obj)"

# Zip up the package for the managed application
./build/managedapp/createAppZip.sh $SCENARIO

# Create the Service Definition
STORAGE_ACC_RESOURCE_GROUP=$MANAGED_APP_STORAGE_RG
STORAGE_ACC_NAME=$MANAGED_APP_STORAGE_NAME
BUILDNUM=$(echo "$MANAGED_APP_VERSION" | awk -F. '{print $3}')
MA_RESOURCE_GROUP=$STORAGE_ACC_RESOURCE_GROUP
MODIFIED_SCENARIO=${SCENARIO//\//-}
MA_DEFINITION_NAME="modm-$MODIFIED_SCENARIO-$MANAGED_APP_VERSION"
MODIFIED_VERSION=${MANAGED_APP_VERSION//./-}
STORAGE_CONTAINER_NAME="modm$MODIFIED_SCENARIO-$MODIFIED_VERSION"

./build/managedapp/createServiceDefinition.sh $STORAGE_ACC_RESOURCE_GROUP $STORAGE_ACC_NAME $STORAGE_CONTAINER_NAME $MA_RESOURCE_GROUP $MA_DEFINITION_NAME

