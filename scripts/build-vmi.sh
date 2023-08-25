#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script delegates to the VMI build script
#   
#   EXAMPLES:
#
#     Calls the build with an explicit env var file of your choice that contains all env vars including for packer
#       ./scripts/build-vmi.sh --env ./obj/.env.local
#
#     Calls the build with required arguments, using default env vars at ./obj/.env.pkrvars
#       ./scripts/build-vmi.sh --version 0.0.301 --resource-group rg-name --image-gallery-name my.sig --location eastus
#
# ===========================================================================================

source ./build/vmi/build.sh
