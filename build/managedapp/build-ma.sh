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

DEPLOYED_IMAGE_REFERENCE="$2"
echo "The deployed image reference is: $DEPLOYED_IMAGE_REFERENCE"
UIDEF_FILE="./build/managedapp/$SCENARIO/createUiDefinition.json"
VIEWDEF_FILE="./build/managedapp/$SCENARIO/viewDefinition.json"
TEMP_FILE="./obj/createUiDefinition.json"
TEMP_VIEWDEF_FILE="./obj/viewDefinition.json" 

echo "The UIDEF_FILE is: $UIDEF_FILE"

# Service Definition Vars
STORAGE_ACC_RESOURCE_GROUP=$MANAGED_APP_STORAGE_RG
STORAGE_ACC_NAME=$MANAGED_APP_STORAGE_NAME
BUILDNUM=$(echo "$MANAGED_APP_VERSION" | awk -F. '{print $3}')
MA_RESOURCE_GROUP=$STORAGE_ACC_RESOURCE_GROUP
MODIFIED_SCENARIO=${SCENARIO//\//-}
MA_DEFINITION_NAME="modm-$MODIFIED_SCENARIO-$MANAGED_APP_VERSION"
MODIFIED_VERSION=${MANAGED_APP_VERSION//./-}
STORAGE_CONTAINER_NAME="modm$MODIFIED_SCENARIO-$MODIFIED_VERSION"
FUNCTION_NAME=$STORAGE_CONTAINER_NAME
FUNCTION_URL="https://$FUNCTION_NAME.azurewebsites.net"

# Assign the Reader role to the Managed Application Service Principal
az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope "$DEPLOYED_IMAGE_REFERENCE"

# Use sed to replace <IMAGE_REFERENCE> with the DEPLOYED_IMAGE_REFERENCE
# The -i option has issues with certain platform implementations of the sed command,
# so we use a temporary file for the output and then overwrite the original file
rm $TEMP_FILE 2> /dev/null

echo "prior to sed command, the UIDEF_FILE contains: $(cat $UIDEF_FILE)"

./build/managedapp/createAzureFunction.sh $STORAGE_ACC_RESOURCE_GROUP $STORAGE_ACC_NAME $STORAGE_CONTAINER_NAME
FUNCTION_BLOB=$(az storage blob url --account-name "$STORAGE_ACC_NAME" --container-name "$STORAGE_CONTAINER_NAME" --name functionapp.zip --output tsv)

sed -e "s|<IMAGE_REFERENCE>|$DEPLOYED_IMAGE_REFERENCE|g" -e "s|<ZIPPED_FUNCTION>|$FUNCTION_BLOB|g" -e "s|<HOSTING_PLAN_NAME>|$FUNCTION_NAME|g" "$UIDEF_FILE" > "$TEMP_FILE"
echo "Before sed command, the VIEWDEF_FILE contains: $(cat $VIEWDEF_FILE)"
sed -e "s|<FUNCTION_APP_NAME>|$FUNCTION_NAME|g" "$VIEWDEF_FILE" > "$TEMP_VIEWDEF_FILE"
echo "After sed command, the TEMP_VIEWDEF_FILE contains: $(cat $TEMP_VIEWDEF_FILE)"
echo "the location of the TEMP_VIEWDEF_FILE is: $(pwd)/$TEMP_VIEWDEF_FILE"

rm ./obj/mainTemplate.json 2> /dev/null
# rm ./obj/viewDefinition.json 2> /dev/null
# rm ./obj/viewDefinition.json 2> /dev/null

cp -f ./build/managedapp/$SCENARIO/mainTemplate.json ./obj/mainTemplate.json
# cp -f ./build/managedapp/$SCENARIO/viewDefinition.json ./obj/viewDefinition.json

echo "The ./obj directory contains: $(ls -la ./obj)"

# Zip up the package for the managed application
./build/managedapp/createAppZip.sh $SCENARIO

./build/managedapp/createServiceDefinition.sh $STORAGE_ACC_RESOURCE_GROUP $STORAGE_ACC_NAME $STORAGE_CONTAINER_NAME $MA_RESOURCE_GROUP $MA_DEFINITION_NAME

