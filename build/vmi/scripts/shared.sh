#!/bin/bash


function pad_string() {
  local value=$1
  local final_length=$2
  
  
  local value_length=${#value}
  local pad_amount=$((final_length-value_length))
  local padding=$( printf "%${pad_amount}s" " " )

  echo "${value}${padding}"
}

function print_packer_variables() {
    local pkr_vars=$(env | grep '^PKR_VAR_')
    while IFS= read -r line ; do
        local key=$(echo "$line" | cut -d "=" -f 1)
        local value=$(echo "$line" | cut -d "=" -f 2)

        local key_display=$(pad_string "${key}:" 30)
        echo "${key_display}${value}"; 
    done <<< "$pkr_vars"
}

function get_build_environment() {
  # Check if running in GitHub Actions environment
  if [ -n "$GITHUB_ACTIONS" ]; then
      echo "GitHub"
  else
      echo "Local"
  fi
}

function export_packer_env_vars_from_file() {
    # make sure we have a vars file before proceeding
    env_pkrvars_file=./obj/.env.pkrvars

    if [ ! -f $env_pkrvars_file ];
    then
        echo "./obj/.env.pkrvars file is required."
        exit 1
    fi

    echo "Exporting Packer variables to environment."

    # export packer env variables so they get picked up
    export $(grep -v '^#' $env_pkrvars_file | xargs)
}

function ensure_shared_image_gallery() {
    # parse options as --option optionValue
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--gallery-name)
        name="$2"
        shift # past argument
        shift # past value
        ;;
        -g|--resource-group)
        rg="$2"
        shift # past argument
        shift # past value
        ;;
        -l|--location)
        location="$2"
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

    echo "Confirming shared image gallery [$name] exists."
    az sig show --resource-group $rg --gallery-name $name --query id --output tsv

    # $? will contain the exit status of the last command
    if [ $? -ne 0 ]; then
        echo "Creating gallery [$name] in [$rg]."
        az sig create --resource-group $rg --gallery-name $name --location $location
    fi

    echo ""
}

function ensure_image_definition() {
    # parse options as --option optionValue
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--subscription)
        subscription_id="$2"
        shift # past argument
        shift # past value
        ;;
        -g|--resource-group)
        resource_group="$2"
        shift # past argument
        shift # past value
        ;;
        -i|--image-gallery-name)
        image_gallery_name="$2"
        shift # past argument
        shift # past value
        ;;
        -o|--image-offer)
        image_offer="$2"
        shift # past argument
        shift # past value
        ;;
        -n|--image-name)
        image_name="$2"
        shift # past argument
        shift # past value
        ;;
        -l|--location)
        location="$2"
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

    echo "Confirming image definition [$image_name] exists in the gallery."

    # ensure image definition
    az sig image-definition show \
        --subscription $subscription_id \
        --resource-group $resource_group \
        --gallery-name $image_gallery_name \
        --gallery-image-definition $image_name \
        --output none

    if [ $? -ne 0 ]; then
    echo "Creating image definition [$image_name] in $image_gallery_name"
    
    az sig image-definition create \
        --subscription $subscription_id \
        --resource-group $resource_group \
        --gallery-name $image_gallery_name \
        --gallery-image-definition $image_name \
        --publisher Microsoft \
        --offer $image_offer \
        --sku $image_offer \
        --os-type Linux \
        --location $location

        if [ $? -ne 0 ]; then
            echo "Image definition creation failed."
            exit 1
        fi
    fi
}

# resource group should exist. do not automatically create it
function demand_resource_group() {
  resource_group=$1
  echo "Checking Resource Group."
  az group show -n $resource_group -o none
  result=$?
  echo ""

  if [ $result -ne 0 ]; then
      echo "ERROR: Resource Group [$resource_group] does not exist. aborting vmi build."
      exit 1
  fi
}

function guard_against_empty() {
    # parse options as --option optionValue
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--value)
        value="$2"
        shift # past argument
        shift # past value
        ;;
        -m|--error-message)
        message="$2"
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

    if [ -z "$value" ]; then
        echo "$message"
        exit 1
    fi
}

function execute_packer() {
    # parse options as --option optionValue
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--image-name)
        name="$2"
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

    # Run the Packer command
    packer_file="./build/vmi/${name}/vmi.pkr.hcl"
    echo "Executing Packer build [$packer_file]"

    packer init $packer_file
    packer build $packer_file
}

export -f pad_string
export -f print_packer_variables
export -f get_build_environment
export -f demand_resource_group
export -f ensure_shared_image_gallery
export -f ensure_image_definition
export -f execute_packer
export -f guard_against_empty