#!/bin/bash

# Dependencies:
# - GitHub CLI - https://cli.github.com/manual/
#     MacOS: brew install gh


function gh_login() {
  repository=microsoft/commercial-marketplace-offer-deploy

  gh auth login -h github.com --web
  gh repo set-default $repository
}

function gh_set_variables() {
  # parse options as --option optionValue
  POSITIONAL_ARGS=()

  while [[ $# -gt 0 ]]; do
  case $1 in
      -e|--env)
      env="$2"
      shift # past argument
      shift # past value
      ;;
      -f|--env-file)
      env_file="$2"
      shift # past argument
      shift # past value
      ;;
      -*|--*)
      echo "Unknown option $1"
      exit 1
      ;;
      *)
      POSITIONAL_ARGS+=("$1") # save positional arg
      shift # past argument
      ;;
  esac
  done
  set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

  echo "Your current auth status with GitHub:"
  gh auth status
  echo ""

  echo "Setting variables for GitHub environment [$env]"
  echo "File: $env_file"

  gh variable set --env $env -f $env_file
}

function gh_set_secrets() {
  # parse options as --option optionValue
  POSITIONAL_ARGS=()

  while [[ $# -gt 0 ]]; do
  case $1 in
      -e|--env)
      env="$2"
      shift # past argument
      shift # past value
      ;;
      -f|--env-file)
      env_file="$2"
      shift # past argument
      shift # past value
      ;;
      -*|--*)
      echo "Unknown option $1"
      exit 1
      ;;
      *)
      POSITIONAL_ARGS+=("$1") # save positional arg
      shift # past argument
      ;;
  esac
  done
  set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

  echo "Your current auth status with GitHub:"
  gh auth status
  echo ""

  echo "Setting secrets for GitHub environment [$env]"
  echo "File: $env_file"

  gh secret set --env $env -f $env_file
}

export -f gh_login
export -f gh_set_variables