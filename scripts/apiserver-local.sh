#!/bin/bash


go build -o ./bin/ ./cmd/apiserver

cp ./configs/.env.development.local ./bin/.env

cd ./bin
./apiserver