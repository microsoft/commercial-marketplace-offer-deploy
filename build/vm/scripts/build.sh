#!/bin/bash

# build 
cd /usr/local/source

sudo docker build ./src -t modm -f ./build/container/Dockerfile.modm  
sudo docker build . -t jenkins -f ./build/container/Dockerfile.jenkins

# next, setup caddy and compose

# write the fqdn to an env file that will be used by caddy for its environment variables
# this will allow the site address to be set to the VM's FQDN
echo VM_FQDN=$VM_FQDN | sudo tee /usr/local/source/build/vm/caddy/.env

sudo docker compose -f /usr/local/source/build/vm/docker-compose.yml -p modm up --force-recreate
# sudo docker compose -f /usr/local/source/build/vm/docker-compose.yml -p modm down