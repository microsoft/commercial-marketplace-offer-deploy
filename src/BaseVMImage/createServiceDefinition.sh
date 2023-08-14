#!/bin/bash

# az group create --name bobjacPackageStorage --location eastus2 

az storage account create \
    --name bobjacmodmscsa29  \
    --resource-group bobjacPackageStorage \
    --location eastus2  \
    --sku Standard_LRS \
    --kind StorageV2

az storage container create \
    --account-name bobjacmodmscsa29 \
    --name appcontainer \
    --auth-mode login \
    --public-access blob

az storage blob upload \
    --account-name bobjacmodmscsa29 \
    --container-name appcontainer \
    --name "app.zip" \
    --file "./app.zip"

blob=$(az storage blob url --account-name bobjacmodmscsa29 --container-name appcontainer --name app.zip --output tsv)
groupid=$(az ad group show --group bobjacmarketplace --query id --output tsv)
roleid=$(az role definition list --name Owner --query [].name --output tsv)

az group create --name appDefinitionGroup36 --location eastus2
az managedapp definition create --name "MODM14" --location "eastus2" --resource-group appDefinitionGroup36 --lock-level ReadOnly --display-name "MODM14" --description "MODM14" --authorizations "$groupid:$roleid" --package-file-uri "$blob"

# This permission was needed to allow the service catalog to work with the image
# az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope /subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/modmImage/providers/Microsoft.Compute/images/bobjacnginx2

