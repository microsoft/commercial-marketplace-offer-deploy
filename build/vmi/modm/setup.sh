#!/bin/bash

# copy files to the correct location

sudo cp /tmp/modm.service /etc/systemd/system/modm.service
sudo cp /tmp/Caddyfile $MODM_HOME/Caddyfile
sudo cp /tmp/Caddyfile $MODM_HOME/docker-compose.yml

# build the docker images

cd $MODM_HOME/source

sudo git checkout $MODM_REPO_BRANCH
sudo git pull

# build final version
sudo docker build ./src -t modm -f ./build/container/Dockerfile.modm  
sudo docker build . -t jenkins -f ./build/container/Dockerfile.jenkins