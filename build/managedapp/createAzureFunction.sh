#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 <storage_account_name> <container_name>"
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

