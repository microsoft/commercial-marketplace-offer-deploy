#!/bin/bash

# This is called by the modm.service to start up the docker compose for MODM

# get the VM instance metadata info by using Azure Instance Metadata Service (IMDS)
# see: https://learn.microsoft.com/en-us/azure/virtual-machines/instance-metadata-service
function get_site_address() {
  metadata_url="http://169.254.169.254/metadata/instance?api-version=2020-09-01"
  curl -H Metadata:true $metadata_url
}

az login --identity

$spID=$(az resource list -n <VM-NAME> --query [*].identity.principalId --out tsv)
echo The managed identity for Azure resources service principal ID is $spID


echo SITE_ADDRESS=$VM_FQDN | sudo tee $MODM_HOME/.env
echo ACME_ACCOUNT_EMAIL=nowhere@nowhere.com  | sudo tee --append $MODM_HOME/.env

# Start up docker compose
# TODO: call entrypoint using docker run instead, attach the .env file
sudo docker compose -f $MODM_HOME/docker-compose.yml -p modm up -d --force-recreate