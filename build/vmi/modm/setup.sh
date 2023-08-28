#!/bin/bash

# build 
cd /usr/local/source

sudo git pull

sudo docker build ./src -t modm -f ./build/container/Dockerfile.modm  
sudo docker build . -t jenkins -f ./build/container/Dockerfile.jenkins
sudo docker build . -t entrypoint-builder -f ./build/container/Dockerfile.modmentry

sudo docker create --name entrypoint-container entrypoint-builder
sudo docker cp entrypoint-container:/app/. /usr/local/bin/
sudo docker rm entrypoint-container