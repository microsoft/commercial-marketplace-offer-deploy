#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script gets the version to use for the next build of a VMI image
#   
#   Prescendence of values
#     It takes two arguments:
#        image_name - the name of the image you want to build
#        resource_group - the resource group where the image should be deployed 
#   File Paths:
#     ALL paths are relative to the root of the repository. This guarantees consistency with relative paths
#     [realpath] is used for converting all relative paths to absolute paths
#
#   Example Usage:
#   next_version=$(get_next_image_version "modmvmi" "modm-dev-vmi")
#
# ===========================================================================================
get_next_image_version() {
    # parse options as --option optionValue
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--gallery-name)
        gallery_name="$2"
        shift # past argument
        shift # past value
        ;;
        -g|--resource-group)
        resource_group="$2"
        shift # past argument
        shift # past value
        ;;
        -i|--image-name)
        image_name="$2"
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

    # get latest version
    latest_version=$(az sig image-version list -g $resource_group -i $image_name -r $gallery_name -o tsv --query "[-1:].name")

    # Increment the patch version
    IFS="." read -r -a version_parts <<< "$latest_version"
    new_patch=$((version_parts[2] + 1))
    next_version="${version_parts[0]}.${version_parts[1]}.$new_patch"

    # Return the formatted next version without image name
    echo "$next_version"
}