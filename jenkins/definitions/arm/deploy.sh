#!/bin/bash

cd $MODM_HOME/installer

# if the Azure client secret is not set, use MSI
if [ -z "$AZURE_CLIENT_SECRET" ]; then
  az login --identity
fi

# special case where echoing this will strip the Jenkins preamble
# from the job output
echo "-----------------"


parameters_file="parameters.json"
template_file=$(cat ./manifest.json | jq -r '.mainTemplate')

resource_group_name=$(cat ./$parameters_file | jq -r '.parameters.resourceGroupName.value')

# Check if resource_group_name is empty
if [ -z "$resource_group_name" ]; then
  echo "Parameter 'resourceGroupName' not found in $parameters_file"
  exit 1
fi

# Deploy using the extracted resource group name
az deployment group create --resource-group $resource_group_name \
    --name deployment1 \
    --template-file $template_file \
    --parameters @$parameters_file