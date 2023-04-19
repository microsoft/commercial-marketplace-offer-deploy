#!/bin/bash

echo "Building MODM image: ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
docker build -t ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} . -f ./build/package/Dockerfile
docker tag ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} ${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest