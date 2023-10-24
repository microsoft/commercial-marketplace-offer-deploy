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
#       ./scripts/build-managedapp.sh -v 0.0.301 \
#           --scenario terraform/simple \
#           --version 0.1.20 \
#           --resource-group modm-dev \
#           --image-id /subscriptions/.../images/modm/versions/0.0.222 \
#           --storage-account modmsvccatalog
# ===========================================================================================

source ./build/managedapp/build.sh
