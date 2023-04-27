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

function run_docker() {
    echo "Building modm container image."
    docker build . -t modm:latest -f ./build/package/Dockerfile --quiet
    
    echo "starting NGROK"
    # start up ngrok and get address
    ngrok http 8080 > /dev/null &
    NGROK_ID=$!
    echo "Pid: $NGROK_ID"
    sleep 2
    export PUBLIC_DOMAIN_NAME=$(curl -s localhost:4040/api/tunnels | jq '.tunnels[0].public_url' -r | sed -E 's/^\s*.*:\/\///g')
    echo "Public Domain: $PUBLIC_DOMAIN_NAME"
    
    docker compose -f ./deployments/docker-compose.standalone.yml up 
}

function kill_ngrok() {
  kill $NGROK_ID
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

  docker)
    trap kill_ngrok EXIT
    run_docker
    ;;
  *)
    echo -n "unknown"
    ;;
esac

