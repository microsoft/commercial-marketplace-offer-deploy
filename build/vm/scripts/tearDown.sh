#!/bin/bash

# Check if all required parameters are provided
if [ $# -ne 3 ]; then
  echo "Usage: $0 <start_version> <end_version> <common_resource_group>"
  exit 1
fi

# Set the range of numbers
start_number=$1
end_number=$2
common_resource_group=$3

# Loop through the range and construct resource group names
for ((number=$start_number; number<=$end_number; number++))
do
    resource_group_name="modm$number"

    # Get a list of associated resource group names that match the pattern
    associated_resource_group_names=$(az group list --query "[?name | starts_with(@, 'rg-modm$number')].name" --output tsv)

    # Loop through the associated resource group names and delete them
    for associated_rg_name in $associated_resource_group_names
    do
        echo "Deleting associated resource group: $associated_rg_name"
        az group delete --name "$associated_rg_name" --yes --no-wait
    done

    # Delete the main resource group
    echo "Deleting resource group: $resource_group_name"
    az group delete --name "$resource_group_name" --yes --no-wait
done
