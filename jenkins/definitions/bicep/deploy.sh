#!/bin/bash

echo "Hello crom deploy.sh"

cd $MODM_HOME/installer

if [ -z "$AZURE_CLIENT_SECRET" ]; then
  az login --identity
  export ARM_USE_MSI=true
  export ARM_CLIENT_ID=$AZURE_CLIENT_ID
  export ARM_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
  export ARM_TENANT_ID=$AZURE_TENANT_ID
else
  # Set Azure credentials from Jenkins bindings
  export ARM_CLIENT_ID=$AZURE_CLIENT_ID
  export ARM_CLIENT_SECRET=$AZURE_CLIENT_SECRET
  export ARM_SUBSCRIPTION_ID=$AZURE_SUBSCRIPTION_ID
  export ARM_TENANT_ID=$AZURE_TENANT_ID
fi

# Extract resource group name from params.json using jq
resourceGroupName=$(jq -r '.parameters.resourceGroupName.value' params.json)

# Check if resourceGroupName is empty
if [ -z "$resourceGroupName" ]; then
  echo "Resource group name not found in params.json"
  exit 1
fi

# Deploy using the extracted resource group name
az deployment group create --resource-group "$resourceGroupName" --name deployment1 --template-file ./main.json --parameters @params.json