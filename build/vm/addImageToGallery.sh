#!/bin/bash

if [ $# -ne 8 ]; then
  echo "Usage: $0 <gallery_resource_group> <gallery_name> <gallery_image_definition> <gallery_image_version> <subscription_id>, <image>, <image_resource_group>, <region"
  exit 1
fi

GALLERY_RESOURCE_GROUP="$1"
GALLERY_NAME="$2"
GALLERY_IMAGE_DEFINITION="$3"
GALLERY_IMAGE_VERSION="$4"
SUBSCRIPTION_ID="$5"
IMAGE="$6"
IMAGE_RESOURCE_GROUP="$7"
REGION="$8"

az sig show --resource-group $GALLERY_RESOURCE_GROUP --gallery-name $GALLERY_NAME --query id --output tsv
# $? will contain the exit status of the last command
if [ $? -eq 0 ]; then
    echo "The gallery exists."
else
    echo "The gallery does not exist."
    az sig create --resource-group "$GALLERY_RESOURCE_GROUP" --gallery-name "$GALLERY_NAME" --location "$REGION"
fi

# Check if the image definition exists
az sig image-definition show --resource-group "$GALLERY_RESOURCE_GROUP" --gallery-name "$GALLERY_NAME" --gallery-image-definition "$GALLERY_IMAGE_DEFINITION" --query id --output tsv

# $? will contain the exit status of the last command
if [ $? -eq 0 ]; then
    echo "The image definition exists."
else
    echo "The image definition does not exist."

    az sig image-definition create \
    --resource-group "$GALLERY_RESOURCE_GROUP" \
    --gallery-name "$GALLERY_NAME" \
    --gallery-image-definition "$GALLERY_IMAGE_DEFINITION" \
    --publisher microsoftgps --offer modm --sku "ubuntu-20.04-custom-app" --os-type Linux
fi

# Check if the image version exists
az sig image-version show --resource-group "$GALLERY_RESOURCE_GROUP" --gallery-name "$GALLERY_NAME" --gallery-image-definition "$GALLERY_IMAGE_DEFINITION" --gallery-image-version "$GALLERY_IMAGE_VERSION" --query id --output tsv

# $? will contain the exit status of the last command
if [ $? -eq 0 ]; then
    echo "The image version exists."
else
    echo "The image version does not exist."

    az sig image-version create \
   --resource-group "$GALLERY_RESOURCE_GROUP" \
   --gallery-name "$GALLERY_NAME" \
   --gallery-image-definition "$GALLERY_IMAGE_DEFINITION" \
   --gallery-image-version "$GALLERY_IMAGE_VERSION" \
   --managed-image "/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$IMAGE_RESOURCE_GROUP/providers/Microsoft.Compute/images/$IMAGE"
fi


