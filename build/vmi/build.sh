#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script defines args/options for the VM Image (VMI) build that is driven by packer
#   
#   DEPENDENCIES:
#     args.sh - This script is required to be run first in order to get all input args set
# 
# ===========================================================================================

source ./build/vmi/args.sh

# verify arg inputs that were processed from args.sh

echo ""
echo "------------------------------------"
echo "ENV FILE        = ${ENV_VARS_FILE}"
echo "VERSION         = ${VERSION}"
echo "RESOURCE GROUP  = ${RESOURCE_GROUP}"
echo "GALLERY NAME    = ${GALLERY_NAME}"
echo "LOCATION        = ${LOCATION}"
echo "CONFIRMED       = ${CONFIRMED} (bypass prompts to proceed)"
echo ""
echo "PACKER VARIABLES:"
echo "  client_id            = ${PKR_VAR_client_id}"
echo "  client_secret        = ***"
echo "  subscription_id      = ${PKR_VAR_subscription_id}"
echo "  tenant_id            = ${PKR_VAR_tenant_id}"
echo "  location             = ${PKR_VAR_location}"
echo "  managed_image_name   = ${PKR_VAR_managed_image_name}"
echo "  resource_group_name  = ${PKR_VAR_resource_group_name}"
echo ""


function ensure_shared_image_gallery() {
  echo "Checking shared image gallery [$GALLERY_NAME]"
  az sig show --resource-group $RESOURCE_GROUP --gallery-name $GALLERY_NAME --query id --output tsv

  # $? will contain the exit status of the last command
  if [ $? -ne 0 ]; then
      echo "Gallery does not exist. Creating gallery [$GALLERY_NAME] in [$RESOURCE_GROUP]"
      az sig create --resource-group $RESOURCE_GROUP --gallery-name $GALLERY_NAME --location $LOCATION
  fi

  echo ""
}

# we're going to have all incoming input variable values come from environment variables, not a var file
function packer_build() {
  echo "Executing Packer build."
  
  packer init ./build/vmi/modm.pkr.hcl
  packer build ./build/vmi/modm.pkr.hcl
}


# -------------------------------------------------------------------------
# main

ensure_shared_image_gallery
packer_build