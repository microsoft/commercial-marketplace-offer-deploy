#!/bin/bash

docker build ./src -t jenkins -f ../build/package/Dockerfile.jenkins  
docker tag modm:latest modmdev.azurecr.io/modm:latest

docker build . -t modm -f ../build/package/Dockerfile.modm
docker tag jenkins:latest modmdev.azurecr.io/jenkins:latest


docker push modmdev.azurecr.io/modm:latest
docker push modmdev.azurecr.io/jenkins:latest
