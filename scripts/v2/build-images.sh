#!/bin/bash

docker build . -t jenkins -f ../build/package/Dockerfile.jenkins  
docker build . -t modm -f ../build/package/Dockerfile.modm


az login