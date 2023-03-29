#!/bin/bash


go build -o ./bin/ ./cmd/operator

cp -n ./configs/.env.development.local ./bin/.env

cd ./bin
./operator