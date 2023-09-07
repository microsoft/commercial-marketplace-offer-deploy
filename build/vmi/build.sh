#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script drives the packer builds
#   
#   imports
    source ./build/vmi/scripts/shared.sh
#
#   Example Call (from root of repository):
#       ./build/vmi/build.sh -v 0.0.301 -n modm-base -o modm-base-ubuntu
# ===========================================================================================

# args
# ------------------------------------------------------------

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
case $1 in
    -n|--image-name)
    IMAGE_NAME="$2"
    shift # past argument
    shift # past value
    ;;
    -v|--image-version)
    IMAGE_VERSION="$2"
    shift # past argument
    shift # past value
    ;;
    -o|--image-offer)
    IMAGE_OFFER="$2"
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

guard_against_empty --value "$IMAGE_NAME" --error-message "Image Name is required. Use --image-name|-n to set."
guard_against_empty --value "$IMAGE_VERSION" --error-message "Image Version is required. Use --image-version|-v to set."
guard_against_empty --value "$IMAGE_VERSION" --error-message "Image Offer is required. Use --image-offer|-o to set."

# Defaults
# ------------------------------------------------------------

RESOURCE_GROUP_NAME="modm-dev-vmi"
BUILD_ENVIRONMENT=$(get_build_environment)


if [ "$BUILD_ENVIRONMENT" = "Local" ]; then
    export_packer_env_vars_from_file
fi

# override default resource group name if packer var declares the resource group
if [ -n "$PKR_VAR_resource_group" ]; then
    RESOURCE_GROUP_NAME=$PKR_VAR_resource_group
fi

export PKR_VAR_image_name="${IMAGE_NAME}"
export PKR_VAR_image_version="${IMAGE_VERSION}"

# DEMANDS
demand_resource_group $RESOURCE_GROUP_NAME

# PREAMBLE
echo ""
echo "Environment Variables "
echo "---------------------------------------------"

echo "BUILD_ENVIRONMENT:        $BUILD_ENVIRONMENT"
echo "RESOURCE_GROUP_NAME:      $RESOURCE_GROUP_NAME"
echo "IMAGE_NAME:               $IMAGE_NAME"
echo "IMAGE_VERSION:            $IMAGE_VERSION"
echo ""
echo "Packer Variables "
echo "---------------------------------------------"
print_packer_variables
echo ""

ensure_shared_image_gallery -g $RESOURCE_GROUP_NAME -n $PKR_VAR_image_gallery_name -l $PKR_VAR_location

ensure_image_definition -s $PKR_VAR_subscription_id -g $RESOURCE_GROUP_NAME -l $PKR_VAR_location \
    --image-name $PKR_VAR_image_name \
    --image-gallery-name $PKR_VAR_image_gallery_name \
    --image-offer $IMAGE_OFFER

execute_packer --image-name $IMAGE_NAME