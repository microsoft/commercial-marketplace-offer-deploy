#!/bin/bash

echo "Building MODM apiserver image: ${CONTAINER_REGISTRY}/${APISERVER_IMAGE_NAME}:${IMAGE_TAG}"
docker build -t ${CONTAINER_REGISTRY}/${APISERVER_IMAGE_NAME}:${IMAGE_TAG} . -f ./build/package/Dockerfile.apiserver
docker tag ${CONTAINER_REGISTRY}/${APISERVER_IMAGE_NAME}:${IMAGE_TAG} ${CONTAINER_REGISTRY}/${APISERVER_IMAGE_NAME}:latest

echo "Building MODM operator image: ${CONTAINER_REGISTRY}/${OPERATOR_IMAGE_NAME}:${IMAGE_TAG}"
docker build -t ${CONTAINER_REGISTRY}/${OPERATOR_IMAGE_NAME}:${IMAGE_TAG} . -f ./build/package/Dockerfile.operator
docker tag ${CONTAINER_REGISTRY}/${OPERATOR_IMAGE_NAME}:${IMAGE_TAG} ${CONTAINER_REGISTRY}/${OPERATOR_IMAGE_NAME}:latest