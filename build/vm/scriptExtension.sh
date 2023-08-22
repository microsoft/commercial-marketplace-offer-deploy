#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Usage: $0 <artifacts location>"
  exit 1
fi

export ARTIFACTS_LOCATION="$1"
export VM_FQDN="$2"

echo "Hello, world from script extension!  The _artifactsLocation is $ARTIFACTS_LOCATION" > scriptExtensionOutput.txt
echo "FQDN = $VM_FQDN"

# next, setup caddy's required .env file. This file is referenced in the docker-compose file for caddy to run properly and compose

# write the fqdn to an env file that will be used by caddy for its environment variables
# this will allow the site address to be set to the VM's FQDN
echo SITE_ADDRESS=$VM_FQDN | sudo tee /usr/local/source/build/vm/caddy/.env
echo ACME_ACCOUNT_EMAIL=nowhere@nowhere.com  | sudo tee --append /usr/local/source/build/vm/caddy/.env

sudo docker compose -f /usr/local/source/build/vm/docker-compose.yml -p modm up -d --force-recreate
# sudo docker compose -f /usr/local/source/build/vm/docker-compose.yml -p modm down