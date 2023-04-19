#!/bin/bash


function run_apiserver() {
    go build -o ./bin/ ./cmd/apiserver
    cp -n ./configs/.env.development.local ./bin/.env

    cd ./bin
    ./apiserver
}

function run_operator() {
    go build -o ./bin/ ./cmd/operator
    cp -n ./configs/.env.development.local ./bin/.env

    cd ./bin
    ./operator
}

process=$1
echo "Target: $process"

case $process in

  operator)
    run_operator
    ;;

  apiserver)
    run_apiserver
    ;;

  *)
    echo -n "unknown"
    ;;
esac