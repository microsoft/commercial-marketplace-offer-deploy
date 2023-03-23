#!/bin/bash


go build -o ./bin/ ./cmd/apiserver

cp -n ./configs/.env.development.local ./bin/.env

cd ./bin
./apiserver