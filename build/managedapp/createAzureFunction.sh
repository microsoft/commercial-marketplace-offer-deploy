#!/bin/bash

if [ $# -ne 3 ]; then
  echo "Usage: $0 <storage_account_resource_group> <storage_account_name> <container_name>"
  exit 1
fi

echo "inside createAzureFunction.sh"

STORAGE_ACC_RESOURCE_GROUP="$1"
STORAGE_ACC_NAME="$2"
STORAGE_CONTAINER_NAME="$3"

original_dir=$(pwd)

cd "./src/Functions" || exit

echo "The current directory is: $(pwd). Publishing the Azure Function."
# Publish the Azure Function
dotnet publish -c Release -o ./publish

echo "ls -la ./publish: $(ls -la ./publish)"

# Zip the output
cd publish
echo "zipping the Azure Function."
zip -r $original_dir/bin/functionapp.zip .
echo "ls -la $original_dir/bin: $(ls -la $original_dir/bin)"

az storage account show --name "$STORAGE_ACC_NAME" --resource-group "$STORAGE_ACC_RESOURCE_GROUP" &>/dev/null
if [ $? -eq 0 ]; then
    echo "Storage account $STORAGE_ACC_NAME exists in resource group $STORAGE_ACC_RESOURCE_GROUP."
else
    echo "Storage account $STORAGE_ACC_NAME does not exist in resource group $STORAGE_ACC_RESOURCE_GROUP."
    echo "Creating storage account $STORAGE_ACC_NAME in resource group $STORAGE_ACC_RESOURCE_GROUP."
    az storage account create \
        --name "$STORAGE_ACC_NAME"  \
        --resource-group "$STORAGE_ACC_RESOURCE_GROUP" \
        --location eastus2  \
        --sku Standard_LRS \
        --kind StorageV2 \
        --allow-blob-public-access true
fi

echo "Creating storage account container $STORAGE_CONTAINER_NAME in storage account $STORAGE_ACC_NAME."
az storage container create \
    --account-name "$STORAGE_ACC_NAME" \
    --name "$STORAGE_CONTAINER_NAME" \
    --auth-mode login \
    --public-access blob

az storage blob upload \
    --account-name $STORAGE_ACC_NAME \
    --container-name $STORAGE_CONTAINER_NAME \
    --type application/zip \
    --name functionapp.zip \
    --type application/zip \
    --file $original_dir/bin/functionapp.zip

