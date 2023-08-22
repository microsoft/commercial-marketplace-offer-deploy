#!/bin/bash

# build 
cd /usr/local/source

sudo docker build ./src -t modm -f ./build/container/Dockerfile.modm  
sudo docker build . -t jenkins -f ./build/container/Dockerfile.jenkins

# next, setup caddy
SITE_ADDRESS=""
sudo docker compose -f /usr/local/source/build/vm/docker-compose.yml -p modm up -d --force-recreate
# sudo docker compose -f /usr/local/source/build/vm/docker-compose.yml -p modm down