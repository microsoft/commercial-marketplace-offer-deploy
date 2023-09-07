#!/bin/bash

# Check if the test_scenario variable is provided
if [ -z "$1" ]; then
  echo "Error: Please provide a test_scenario as an argument."
  exit 1
fi

test_scenario=$1

# Check if the directory with the test_scenario name exists
if [ ! -d "./v2/$test_scenario" ]; then
  echo "Error: Directory '$test_scenario' does not exist."
  exit 1
fi

# Zip the contents of the directory
cd "./v2/$test_scenario"
zip -r ../../content.zip .

# Check if the storage account exists
az storage account show --name modmdrop --resource-group modm-dev &> /dev/null
if [ $? -ne 0 ]; then
  # Storage account does not exist, create it
  az storage account create --name modmdrop --resource-group modm-dev --sku Standard_LRS --kind StorageV2 --location eastus
fi

# Check if the container exists
az storage container show --name artifacts --account-name modmdrop &> /dev/null
if [ $? -ne 0 ]; then
  # Container does not exist, create it
  az storage container create --name artifacts --account-name modmdrop &> /dev/null
fi

# Generate a SAS token for the "artifacts" container with read permissions
sas_token="sp=racwdli&st=2023-09-07T17:42:54Z&se=2024-07-01T01:42:54Z&spr=https&sv=2022-11-02&sr=c&sig=RdEDrbswin7T84hQtlI1eWf4AY9NMwsfG56laN1sLpg%3D"
echo "SAS token: $sas_token"

# Upload the contents.zip file to the "artifacts" container
az storage blob upload --container-name artifacts --file ../../content.zip --name content.zip --account-name modmdrop --sas-token $sas_token

echo "Scenario '$test_scenario' zipped and uploaded to modmdrop storage account."

