#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script is executed during the packer VMI build. It will be checked out via git
#   running on the VM. Using the version value, it will get the correct version of the script
#   by using the version as the tagged version and then git checkout that version
# 
# ===========================================================================================

# parse args
POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -v|--version)
      VERSION="$2"
      shift # past argument
      shift # past value
      ;;
    *)
      POSITIONAL_ARGS+=("$1") # save positional arg
      shift # past argument
      ;;
  esac
done
set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

# parse args

echo ""
echo "------------------------------------"
echo "MODM Version: $VERSION"

# # move 
# cd $MODM_HOME/source
# git checkout tags/v1.0 -b v1.0-branch


git checkout tags/v1.0 -b v1.0-branch


# sudo git pull

# sudo docker build ./src -t modm -f ./build/container/Dockerfile.modm  
# sudo docker build . -t jenkins -f ./build/container/Dockerfile.jenkins