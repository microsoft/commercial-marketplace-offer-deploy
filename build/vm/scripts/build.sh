#!/bin/bash

# build 
cd /usr/local/source

docker build ./src -t modm -f ./build/container/Dockerfile.modm  
docker build . -t jenkins -f ./build/container/Dockerfile.jenkins

# next, setup caddy
SITE_ADDRESS=""
sudo docker compose -f ./docker-compose.yml -p modm up