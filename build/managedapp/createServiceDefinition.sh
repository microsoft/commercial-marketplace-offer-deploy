#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 5 ]; then
  echo "Usage: $0 <storage_account_resource_group> <storage_account_name> <container_name> <managed_app_resource_group> <managed_app_definition_name>"
  exit 1
fi
STORAGE_ACC_RESOURCE_GROUP="$1"
STORAGE_ACC_NAME="$2"
STORAGE_CONTAINER_NAME="$3"
MA_RESOURCE_GROUP="$4"
MA_DEFINITION_NAME="$5"

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
    --account-name "$STORAGE_ACC_NAME" \
    --container-name "$STORAGE_CONTAINER_NAME" \
    --name "app.zip" \
    --file "../../bin/app.zip"

blob=$(az storage blob url --account-name "$STORAGE_ACC_NAME" --container-name "$STORAGE_CONTAINER_NAME" --name app.zip --output tsv)
groupid=$(az ad group show --group "Managed Application Tests" --query id --output tsv)
roleid=$(az role definition list --name Owner --query [].name --output tsv)

# az group create --name "$MA_RESOURCE_GROUP" --location eastus2

az managedapp definition create --name "$MA_DEFINITION_NAME" \
    --location "eastus2" \
    --resource-group "$MA_RESOURCE_GROUP" \
    --lock-level ReadOnly \
    --display-name "$MA_DEFINITION_NAME" \
    --description "$MA_DEFINITION_NAME" \
    --authorizations "$groupid:$roleid" \
    --package-file-uri "$blob"
