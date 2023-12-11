#!/bin/bash

# ===========================================================================================
#
#   DESCRIPTION:
#   This script drives the managed application build
#   
    # imports
    source ./scripts/cli-setup.sh
#
#   Example Call (from root of repo):
#
#       ./build/managedapp/build.sh -v 0.0.301 \
#           --scenario terraform/simple \
#           --version 0.1.20 \
#           --resource-group modm-dev \
#           --image-id /subscriptions/.../images/modm/versions/0.0.222 \
#           --storage-account modmsvccatalog

# ===========================================================================================


# shared functions
# ------------------------------------------------------------

function getExpiry() {
    local expiry
    if command -v gdate &> /dev/null; then
        expiry=$(gdate -d "+730 days" '+%Y-%m-%dT%H:%MZ')
    else
        expiry=$(date -v +730d '+%Y-%m-%dT%H:%MZ')
    fi
    echo $expiry
}

function createApplicationPackage() {
    echo "Creating application package."

    # build resources
    info=$(cat $SCENARIO_PATH/manifest.json | jq .)
    echo $info | jq .

    echo "Creating resources tarball."
    echo ""

    modm util create-resources-archive -t ./templates -f src/ClientApp/ClientApp.csproj -o ./bin
    resources_file=./bin/resources.tar.gz

    # Extract mainTemplate, default to "main.tf" if not present
    main_template=$(jq -r '.mainTemplate // "main.tf"' $SCENARIO_PATH/manifest.json)
    

    # build application package, e.g. app.zip
    modm package build \
        --name "$(echo $info | jq .offer.name -r)" \
        --description "$(echo $info | jq .offer.name -r)" \
        --resources-file $resources_file \
        --vmi-reference-id $IMAGE_ID \
        --main-template $SCENARIO_PATH/templates/$main_template \
        --create-ui-definition $SCENARIO_PATH/createUiDefinition.json \
        --out-dir ./bin

    export PACKAGE_FILE=./bin/app.zip

    # Assign the Reader role to the image for the Managed Application Service Principal
    echo "assigning reader role to the VMI"
    az role assignment create --assignee c3551f1c-671e-4495-b9aa-8d4adcd62976 --role acdd72a7-3385-48ef-bd42-f606fba81ae7 --scope $IMAGE_ID -o none
}

function createServiceDefinition() {
    echo "Ensuring storage account [$STORAGE_ACCOUNT_NAME]."
    not_exists=$(az storage account check-name -n $STORAGE_ACCOUNT_NAME -o tsv --query "nameAvailable")

    if [ "$not_exists" = "true" ]; then
        echo "Creating storage account in [$RESOURCE_GROUP]."
        az storage account create \
            --name "$STORAGE_ACCOUNT_NAME"  \
            --resource-group "$RESOURCE_GROUP" \
            --location eastus2  \
            --sku Standard_LRS \
            --kind StorageV2 \
            --allow-blob-public-access true
    fi

    echo "Creating storage account container $STORAGE_CONTAINER_NAME."
    exists=$(az storage container exists --account-name $STORAGE_ACCOUNT_NAME --name $STORAGE_CONTAINER_NAME --auth-mode login --output tsv --query exists)

    if [ "$exists" = "false" ]; then
        az storage container create \
            --account-name "$STORAGE_ACCOUNT_NAME" \
            --name "$STORAGE_CONTAINER_NAME" \
            --auth-mode login \
            --public-access blob
    fi

    az storage blob upload \
        --account-name "$STORAGE_ACCOUNT_NAME" \
        --container-name "$STORAGE_CONTAINER_NAME" \
        --name "app.zip" \
        --file $PACKAGE_FILE \
        --overwrite

    storage_account_connection_string=$(az storage account show-connection-string \
        --name $STORAGE_ACCOUNT_NAME \
        --resource-group $RESOURCE_GROUP \
        --output tsv \
        --query connectionString)
    
    echo "storage_account_connection_string: $storage_account_connection_string"
  
    expiry=$(getExpiry)
    echo "expiry: $expiry"


    echo "Generating SAS token for app.zip with Storage Account Connection String: $storage_account_connection_string expiry: $expiry storage account name: $STORAGE_ACCOUNT_NAME, storage container name: $STORAGE_CONTAINER_NAME"

    blob=$(az storage blob generate-sas \
        --account-name $STORAGE_ACCOUNT_NAME \
        --container-name $STORAGE_CONTAINER_NAME \
        --connection-string $storage_account_connection_string \
        --name app.zip \
        --permissions r \
        --expiry $expiry \
        --https-only \
        --full-uri \
        --output tsv)
    
    
    roleid=$(az role definition list --name Owner --query [].name --output tsv)
    groupid="d391271a-216a-49e1-a36e-c24b2c619f14"

    echo "Application Package Information:"
    echo "Blob              = $blob"
    echo "GroupId           = $groupid"
    echo "RoleId            = $roleid"
    echo "Authorizations    = $groupid:$roleid"
    echo ""

    echo "Creating managed app definition [$NAME]."

    az managedapp definition create --name "$NAME" \
        --location "eastus2" \
        --resource-group "$RESOURCE_GROUP" \
        --display-name "$NAME" \
        --lock-level ReadOnly \
        --description "$NAME" \
        --authorizations "$groupid:$roleid" \
        --package-file-uri "$blob"
}



# args
# ------------------------------------------------------------

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
case $1 in
    -g|--resource-group)
    RESOURCE_GROUP="$2" #the resource group where the managed application is located
    shift # past argument
    shift # past value
    ;;
    -v|--version)
    VERSION="$2"
    shift # past argument
    shift # past value
    ;;
    -i|--image-id)
    IMAGE_ID="$2"
    shift # past argument
    shift # past value
    ;;
    -s|--scenario) # in the format of {deploymentType}/{name}
    SCENARIO="$2"
    SCENARIO=${SCENARIO//$'\n'/}
    SCENARIO=${SCENARIO//$'\r'/}
    shift # past argument
    shift # past value
    ;;
    -a|--storage-account)
    STORAGE_ACCOUNT_NAME="$2" # the storage account to store the service catalog app
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


# Variables based on parameters
echo ""
SCENARIO_NAME=${SCENARIO//\//-}
SCENARIO_PATH=./build/managedapp/$SCENARIO
NAME="modm-$SCENARIO_NAME-$VERSION"
STORAGE_CONTAINER_NAME="$SCENARIO_NAME-${VERSION//./}"


echo "SCENARIO: $SCENARIO"
echo "----------------------------------------------------------------"
echo "  VERSION                         = $VERSION"
echo "  NAME                            = $NAME"
echo "  RESOURCE_GROUP                  = $RESOURCE_GROUP"
echo "  IMAGE_ID                        = $IMAGE_ID"
echo "  SCENARIO_PATH                   = $SCENARIO_PATH"
echo "  STORAGE_ACCOUNT_NAME            = $STORAGE_ACCOUNT_NAME"
echo "  STORAGE_CONTAINER_NAME          = $STORAGE_CONTAINER_NAME"


# main

createApplicationPackage
createServiceDefinition

echo "done."