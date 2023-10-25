#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script delegates to the VMI build script
#   
#   EXAMPLES:
#
#     Calls the build with required arguments, using default env vars at ./obj/.env.pkrvars
#     when called locally
#
#       ./scripts/build-managedapp.sh \
#           --scenario terraform/simple \
#           --version 0.1.134 \
#           --resource-group modm-dev \
#           --image-id /subscriptions/.../images/modm/versions/0.1.134 \
#           --storage-account modmsvccatalog
# ===========================================================================================

source ./build/managedapp/build.sh
