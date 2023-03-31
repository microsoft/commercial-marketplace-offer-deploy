#!/bin/bash

# ingress agent - builds the NGROK ingress agent
#
# example usage: ./ingressagent app=operator -port=8081
#
go build -o ./bin/ ./tools/ingressagent
cp -n ./configs/.env.development.local ./bin/.env
echo "Ingress Agent built"