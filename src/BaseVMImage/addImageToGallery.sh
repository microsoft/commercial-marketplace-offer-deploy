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

az sig image-version create \
   --resource-group "$GALLERY_RESOURCE_GROUP" \
   --gallery-name "$GALLERY_NAME" \
   --gallery-image-definition "$GALLERY_IMAGE_DEFINITION" \
   --gallery-image-version "$GALLERY_IMAGE_VERSION" \
   --managed-image "/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$IMAGE_RESOURCE_GROUP/providers/Microsoft.Compute/images/$IMAGE"

