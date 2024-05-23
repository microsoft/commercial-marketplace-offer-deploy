#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script automates the export of an image from an Azure Image Gallery to a Blob Storage Account.
#
#   Example Call (from root of repository):
#       ./export-vmi.sh -g modm-image-export -l eastus -n bobjactestimage -i /subscriptions/.../images/modm/versions/2.0.208 -a modmimageexport -c images -e 2024-01-08T00:00:00Z
# ===========================================================================================
#   imports
source ./build/vmi/scripts/shared.sh
# args
# ------------------------------------------------------------

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
case $1 in
    -g|--resource-group)
    RESOURCE_GROUP="$2"
    shift # past argument
    shift # past value
    ;;
    -l|--location)
    LOCATION="$2"
    shift # past argument
    shift # past value
    ;;
    -n|--image-name)
    IMAGE_NAME="$2"
    shift # past argument
    shift # past value
    ;;
    -i|--image-id)
    IMAGE_ID="$2"
    shift # past argument
    shift # past value
    ;;
    -a|--account-name)
    STORAGE_ACCOUNT_NAME="$2"
    shift # past argument
    shift # past value
    ;;
    -c|--container-name)
    STORAGE_CONTAINER_NAME="$2"
    shift # past argument
    shift # past value
    ;;
    -e|--expiry)
    EXPIRY_DATE="$2"
    shift # past argument
    shift # past value
    ;;
    -*|--*)
    echo "Unknown option $1"
    exit 1
    ;;
    *)
    POSITIONAL_ARGS+=("$1") # save positional arg
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

# Validate required parameters
guard_against_empty --value "$RESOURCE_GROUP" --error-message "Resource Group is required. Use --resource-group|-g to set."
guard_against_empty --value "$LOCATION" --error-message "Location is required. Use --location|-l to set."
guard_against_empty --value "$IMAGE_NAME" --error-message "Image Name is required. Use --image-name|-n to set."
guard_against_empty --value "$IMAGE_ID" --error-message "Image ID is required. Use --image-id|-i to set."
guard_against_empty --value "$STORAGE_ACCOUNT_NAME" --error-message "Storage Account Name is required. Use --account-name|-a to set."
guard_against_empty --value "$STORAGE_CONTAINER_NAME" --error-message "Storage Container Name is required. Use --container-name|-c to set."
guard_against_empty --value "$EXPIRY_DATE" --error-message "Expiry Date is required. Use --expiry|-e to set."

# Main function
function export_image() {
    echo "Creating disk from image."
    az disk create --resource-group $RESOURCE_GROUP --location $LOCATION --name $IMAGE_NAME --gallery-image-reference $IMAGE_ID

    echo "Granting access to disk."
    access_url=$(az disk grant-access --resource-group $RESOURCE_GROUP --name $IMAGE_NAME --duration-in-seconds 36000 --access-level Read --query accessSas -o tsv)

    echo "Generating SAS for storage container."
    sas_token=$(az storage container generate-sas --account-name $STORAGE_ACCOUNT_NAME --name $STORAGE_CONTAINER_NAME --permissions acw --expiry $EXPIRY_DATE --output tsv)
    destination_url="https://$STORAGE_ACCOUNT_NAME.blob.core.windows.net/$STORAGE_CONTAINER_NAME/$IMAGE_NAME.vhd?$sas_token"

    echo "Copying disk to blob storage."
    azcopy copy "$access_url" "$destination_url" --blob-type PageBlob
}

# Execute main function
export_image

echo "done."
