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
    local image_name="$1"
    local resource_group="$2"
    local latest_version="0.0.0"

    # Get a list of image versions in the resource group
    versions=$(az image list --resource-group "$resource_group" --query "[?starts_with(name, '$image_name-')].name" --output tsv)

    # If no versions found, return "not found" error
    if [ -z "$versions" ]; then
        echo "Image $image_name not found in resource group $resource_group"
        exit 1
    fi

    # Loop through versions to find the latest
    for version in $versions; do
        version_number=$(echo "$version" | awk -F- '{print $NF}' | cut -d'.' -f3)
        if [ "$version_number" -gt "$(echo "$latest_version" | cut -d'.' -f3)" ]; then
            latest_version=$(echo "$version" | awk -F- '{print $NF}')
        fi
    done

    # Increment the patch version
    IFS='.' read -r -a version_parts <<< "$latest_version"
    new_patch=$((version_parts[2] + 1))
    next_version="${version_parts[0]}.${version_parts[1]}.$new_patch"

    # Return the formatted next version without image name
    echo "$next_version"
}

# Example usage:
next_version=$(get_next_image_version "modmvmi" "modm-dev-vmi")
if [ "$next_version" != "not found" ]; then
    echo "Next version for modmvmi: $next_version"
else
    echo "Error: $next_version"
fi

next_version=$(get_next_image_version "modmvmi-base" "modm-dev-vmi")
if [ "$next_version" != "not found" ]; then
    echo "Next version for modmvmi-base: $next_version"
else
    echo "Error: $next_version"
fi

next_version=$(get_next_image_version "hark" "modm-dev-vmi")
if [ "$next_version" != "not found" ]; then
    echo "Next version for hark: $next_version"
else
    echo "Error: $next_version"
fi



