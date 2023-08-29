#!/bin/bash

MODM_VERSION="0.0.0"
source ./build/vmi/scripts/nextversion.sh

if [ $# -ne 1 ]; then
    MODM_VERSION=$(get_next_image_version "modmvmi" "modm-dev-vmi")
else
    MODM_VERSION="$1"
fi

echo "Building modm version $MODM_VERSION"

# Check if running in GitHub Actions environment
if [ -n "$GITHUB_ACTIONS" ]; then
    echo "Running in GitHub Actions environment"
    # You don't need to check for the .env.pkrvars file or export variables here
else
    echo "Running locally"
    # make sure we have a vars file before proceeding
    env_pkrvars_file=./obj/.env.pkrvars

    if [ ! -f $env_pkrvars_file ];
    then
        echo "./obj/.env.pkrvars file is required."
        exit 1
    else
        echo "Packer variables env var file present."
    fi

    # export packer env variables so they get picked up
    export $(grep -v '^#' $env_pkrvars_file | xargs)
fi


export PKR_VAR_sig_image_version_modm=${MODM_VERSION}
export PKR_VAR_managed_image_name_modm=modmvmi-${MODM_VERSION}

# Run the Packer command
packer init ./build/vmi/modm/modm.pkr.hcl
packer build ./build/vmi/modm/modm.pkr.hcl 

