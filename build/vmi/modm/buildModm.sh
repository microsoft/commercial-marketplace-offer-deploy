#!/bin/bash

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

MODM_VERSION="$1"
export PKR_VAR_sig_image_version_modm=${MODM_VERSION}
export PKR_VAR_managed_image_name_modm=modm-${MODM_VERSION}

# Run the Packer command
packer init ./build/vmi/modm/modm.pkr.hcl
packer build ./build/vmi/modm/modm.pkr.hcl 

