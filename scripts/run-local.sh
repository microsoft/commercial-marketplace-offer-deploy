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

function run_modm() {
    arg_build=$1
    
    if [ "$arg_build" = "build" ]  || [ "$arg_build" = "true" ]; then
      build_image modm ./build/package/Dockerfile
    fi

    start_ngrok_background

    compose_file=./deployments/docker-compose.standalone.yml
    docker compose -f $compose_file up 
}

function run_testharness() {
    arg_build=$1

    # the paths are going to be relative to the ./tools dir but the execution is always from the root of the repo dir
    if [ "$arg_build" = "build" ]  || [ "$arg_build" = "true" ]; then
      pushd ../
      build_image modm ./build/package/Dockerfile
      build_image testharness ./build/package/Dockerfile.testharness
      popd
    fi

    start_ngrok_background

    compose_file=../deployments/docker-compose.testharness.yml
    docker compose -f $compose_file up 
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

function build_image() {
  image_name=$1
  dockerfile=$2
  echo "Building $image_name container image."
  docker build . -t $image_name:latest -f $dockerfile --quiet
  docker_build_result=$?

  if [ $docker_build_result -gt 0 ]; then
    echo "docker image build failed."
    echo "exiting."
    exit 1
  fi
}

function start_ngrok_background() {
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
  export MODM_PUBLIC_BASE_URL=$(curl -s localhost:4040/api/tunnels | jq '.tunnels[] | select(.proto == "https") | .public_url' -r)
  echo "NGROK URL:  $MODM_PUBLIC_BASE_URL"
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

  modm)
    trap run_docker_cleanup EXIT
    run_modm $2
    ;;
  testharness)
    trap run_docker_cleanup EXIT
    run_testharness $2
    ;;
  *)
    echo -n "unknown"
    ;;
esac