#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 <artifacts location> <fully qualified domain name>"
  exit 1
fi

export MODM_HOME=/usr/local/modm

export ARTIFACTS_LOCATION="$1"
export VM_FQDN="$2"

echo "Hello, world from script extension!  The _artifactsLocation is $ARTIFACTS_LOCATION" > scriptExtensionOutput.txt
echo "FQDN = $VM_FQDN"

# next, setup caddy's required .env file. This file is referenced in the docker-compose file for caddy to run properly and compose

# write the fqdn to an env file that will be used by caddy and modm for its environment variables
# this will allow the site address to be set to the VM's FQDN
echo SITE_ADDRESS=$VM_FQDN | sudo tee $MODM_HOME/.env
echo ACME_ACCOUNT_EMAIL=nowhere@nowhere.com  | sudo tee --append $MODM_HOME/.env

# Start up docker compose
# TODO: call entrypoint using docker run instead, attach the .env file
sudo docker compose -f $MODM_HOME/docker-compose.yml -p modm up -d --force-recreate