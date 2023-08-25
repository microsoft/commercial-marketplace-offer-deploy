#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script defines args/options for the VM Image (VMI) build that is driven by packer
#   
#   Prescendence of values
#     if an env file is set through --env or via default env location in ./obj, will override any value passed as args to the script
#   File Paths:
#     ALL paths are relative to the root of the repository. This guarantees consistency with relative paths
#     [realpath] is used for converting all relative paths to absolute paths
#
#     MAC OSX: brew install coreutils
#
# ===========================================================================================


# this will be used so that we get a value or set a default value
function getValueOrDefault() {
  if [ -n "$1" ]; then
    echo $1
  else 
    echo $2
  fi
}

function require() {
  if [ -z "$1" ]; then
    echo "$2 is required."
    if [ -n "$3" ]; then
      echo "  " + $3
    fi
    exit 1
  fi
}

function handle_env_var_file() {
  # variables for packer and the script
  # make sure we have a vars file before proceeding
  default_env_pkrvars_file=./obj/.env.pkrvars

  if [ "$ENV_VARS_FILE" = "" ]; then
    if [ ! -f "$default_env_pkrvars_file" ]; then
        echo "environment vars file is required. Please specify a value for --env or place a '.env.pkrvars' in ./obj at the root of the repository"
        exit 1
    else
      ENV_VARS_FILE=$default_env_pkrvars_file
    fi
  fi

  ENV_VARS_FILE=$(realpath $ENV_VARS_FILE)

  # no export all env vars from file so that the explicitly set args can override
  export $(grep -v '^#' $ENV_VARS_FILE | xargs)
}


POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -e|--env)
      ENV_VARS_FILE="$2"
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

handle_env_var_file


# with env var out of the way,
# parse options as --option optionValue
POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
      -v|--version)
      VERSION="$2"
      shift # past argument
      shift # past value
      ;;
    -g|--resource-group)
      RESOURCE_GROUP="$2"
      shift # past argument
      shift # past value
      ;;
    -i|--image-gallery-name)
      GALLERY_NAME=$(getValueOrDefault $2 $GALLERY_NAME)
      shift # past argument
      shift # past value
      ;;
    -l|--location)
      LOCATION="$2"
      shift # past argument
      shift # past value
      ;;
    --yes)
      CONFIRMED=true
      shift # past argument
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


# confirmed
if [ !$CONFIRMED ]; then
  CONFIRMED=false
fi

require "$VERSION" "Version" "use [-v | --version]"
require "$RESOURCE_GROUP" "Resource Group" "use [-g | --resource-group]"
require "$GALLERY_NAME" "Gallery Name" "use [-i | --image-gallery-name]"



# resource group should exist. do not automatically create it
function resource_group_must_exist() {
  echo "Checking Resource Group."
  az group show -n $RESOURCE_GROUP -o json
  result=$?
  echo ""

  if [ $result -ne 0 ]; then
      echo "ERROR: Resource Group [$RESOURCE_GROUP] does not exist. aborting vmi build."
      exit 1
  fi
}

function set_location_to_resource_groups_location_if_not_set() {
  if [ "$LOCATION" = "" ]; then
    LOCATION=$(az group show -n $RESOURCE_GROUP -o tsv --query location)
    echo "Location arg not set. defaulting to resource group's location [$LOCATION]."
  fi
}

function set_packer_vars() {
  export PKR_VAR_location=$LOCATION
  export PKR_VAR_resource_group_name=$RESOURCE_GROUP
  export PKR_VAR_managed_image_name=ubuntu-modm-$VERSION
}


resource_group_must_exist
set_location_to_resource_groups_location_if_not_set
set_packer_vars