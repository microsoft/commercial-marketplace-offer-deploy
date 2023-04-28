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
    arg=$1

    if [ "$arg" = "build" ]  || [ "$arg" = "true" ]; then
      echo "Building modm container image."
      docker build . -t modm:latest -f ./build/package/Dockerfile --quiet
      docker_build_result=$?

      if [ $docker_build_result -gt 0 ]; then
        echo "docker image build failed."
        echo "exiting."
        exit 1
      fi
    fi

    # start up ngrok and get address
    ngrok http 8080 > /dev/null &
    ngrok_start_result=$?
    export NGROK_ID=$!

    if [ $ngrok_start_result -gt 0 ]; then
      echo "NGROK failed to start."
      echo "exiting."
      exit 1
    fi

    echo "NGROK started: $NGROK_ID"
    sleep 2
    export PUBLIC_DOMAIN_NAME=$(curl -s localhost:4040/api/tunnels | jq '.tunnels[0].public_url' -r | sed -E 's/^\s*.*:\/\///g')
    echo "NGROK domain:  $PUBLIC_DOMAIN_NAME \n"

    docker compose -f ./deployments/docker-compose.standalone.yml up 
}

function run_docker_cleanup() {
  echo "Shutdown cleanup..."
  # make sure ngrok is killed
  echo "  killing ngrok"
  kill $NGROK_ID 2>/dev/null

  echo "  removing ready file in ~/tmp"
  rm ~/tmp/ready 2>/dev/null
  echo ""
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
    trap run_docker_cleanup EXIT
    run_docker $2
    ;;
  *)
    echo -n "unknown"
    ;;
esac