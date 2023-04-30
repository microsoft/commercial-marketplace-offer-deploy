#!/bin/bash

# ingress agent - builds the NGROK ingress agent
#
# example usage: ./ingressagent app=operator -port=8081
#
go build -o ./bin/ ./tools/ingressagent
cp -n ./configs/.env.development.local ./bin/.env
echo "Ingress Agent built"

# ingress agent - builds the NGROK ingress agent
#
# example usage: ./testharness
#
go build -o ./bin/ ./tools/testharness

mkdir -p ./bin/template/
cp -n ./test/testdata/taggeddeployment/mainTemplateBicep.json ./bin/template/mainTemplateBicep.json
cp -n ./test/testdata/taggeddeployment/parametersBicep.json ./bin/template/parametersBicep.json

echo "Test Harness built"