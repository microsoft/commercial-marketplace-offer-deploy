#!/bin/bash

# ==========================================================
# Run setup locally using docker compose
# Locally, we are not using Caddy to proxy, so access both modm and jenkins via their configured ports on localhost
#
#
# Example:
#   ./scripts/run-local.sh --build
# ==========================================================

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -b|--build)
      BUILD_IMAGES=$2
      shift # past argument
      shift # past value
      ;;
    *)
      POSITIONAL_ARGS+=("$1") # save positional arg
      shift # past argument
      ;;
  esac
done
set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters


function build_images() {
  if [ $BUILD_IMAGES ]; then
    echo "Building images."
    docker build . -t jenkins -f ./build/container/Dockerfile.jenkins     
    docker build ./src -t modm -f ./build/container/Dockerfile.modm   
  else 
    echo "Skipping image builds."
  fi
}

function run_local() {

  # default location of MODM_HOME for local development
  export MODM_HOME=~/.modm
  build_images
  docker compose -p modm up 
}


run_local