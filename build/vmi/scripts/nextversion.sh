#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script returns the next version to use for the build of a given VMI image
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

get_next_image_version() {
    local image_name="$1"
    local resource_group="$2"
    local latest_version="0.0.0"

    # Get a list of image versions in the resource group
    versions=$(az image list --resource-group "$resource_group" --query "[?starts_with(name, '$image_name-')].name" --output tsv)

    # Loop through versions to find the latest
    for version in $versions; do
        version_number=$(echo "$version" | cut -d'-' -f2- | cut -d'.' -f3)
        if [ "$version_number" -gt "$(echo "$latest_version" | cut -d'.' -f3)" ]; then
            latest_version=$(echo "$version" | cut -d'-' -f2-)
        fi
    done

    # Increment the patch version
    IFS='.' read -r -a version_parts <<< "$latest_version"
    new_patch=$((version_parts[2] + 1))
    next_version="${version_parts[0]}.${version_parts[1]}.$new_patch"

    # Return the formatted next version
    echo "$image_name-$next_version"
}

# Example usage:
next_version=$(get_next_image_version "modmvmi" "modm-dev-vmi")
echo "Next version for modmvmi: $next_version"

next_version=$(get_next_image_version "modmvmi-base" "modm-dev-vmi")
echo "Next version for modmvmi-base: $next_version"

next_version=$(get_next_image_version "hark" "modm-dev-vmi")
echo "Next version for hark: $next_version"