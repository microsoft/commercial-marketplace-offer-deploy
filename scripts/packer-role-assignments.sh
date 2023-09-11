#!/bin/bash

# parse options as --option optionValue
POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
case $1 in
    -a|--assignee)
    packer_client_id="$2"
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

vmi_resource_group=modm-dev-vmi
build_resource_group=modm-dev-packer
shared_image_gallery=modm.dev.sig

image_gallery_id=$(az sig show -r $shared_image_gallery -g $vmi_resource_group --query id --output tsv)
build_resource_group_id=$(az group show -n $vmi_resource_group --query id --output tsv)

contributor_role_id=b24988ac-6180-42a0-ab88-20f7382dd24c
reader_role_id=acdd72a7-3385-48ef-bd42-f606fba81ae7


az role assignment create --assignee $packer_client_id --role $reader_role_id --scope $image_gallery_id
az role assignment create --assignee $packer_client_id --role $contributor_role_id --scope $build_resource_group_id