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
#       ./scripts/build-vmi.sh --version 0.0.1 --image-name modm --image-offer modm-ubuntu
#
# ===========================================================================================

source ./build/vmi/build.sh
