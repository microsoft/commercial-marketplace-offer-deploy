#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 <artifacts location> <fully qualified domain name>"
  exit 1
fi

export MODM_HOME=/usr/local/modm

export ARTIFACTS_LOCATION="$1"
export VM_FQDN="$2"

echo "$ARTIFACTS_LOCATION" > $MODM_HOME/artifacts.uri
echo "FQDN = $VM_FQDN"

# next, setup caddy's required .env file. This file is referenced in the docker-compose file for caddy to run properly and compose

# write the fqdn to an env file that will be used by caddy and modm for its environment variables
# this will allow the site address to be set to the VM's FQDN
echo SITE_ADDRESS=$VM_FQDN | sudo tee $MODM_HOME/.env
echo ACME_ACCOUNT_EMAIL=nowhere@nowhere.com  | sudo tee --append $MODM_HOME/.env
