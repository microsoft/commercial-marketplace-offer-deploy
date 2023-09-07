#!/bin/bash


function pad_string() {
  local value=$1
  local final_length=$2
  
  
  local value_length=${#value}
  local pad_amount=$((final_length-value_length))
  local padding=$( printf "%${pad_amount}s" " " )

  echo "${value}${padding}"
}

function print_packer_variables() {
    local pkr_vars=$(env | grep '^PKR_VAR_')
    while IFS= read -r line ; do
        local key=$(echo "$line" | cut -d "=" -f 1)
        local value=$(echo "$line" | cut -d "=" -f 2)

        local key_display=$(pad_string "${key}:" 30)
        echo "${key_display}${value}"; 
    done <<< "$pkr_vars"
}

function get_build_environment() {
  # Check if running in GitHub Actions environment
  if [ -n "$GITHUB_ACTIONS" ]; then
      echo "GitHub"
  else
      echo "Local"
  fi
}

function export_packer_env_vars_from_file() {
      # make sure we have a vars file before proceeding
    env_pkrvars_file=./obj/.env.pkrvars

    if [ ! -f $env_pkrvars_file ];
    then
        echo "./obj/.env.pkrvars file is required."
        exit 1
    fi

    echo "Exporting Packer variables to environment."

    # export packer env variables so they get picked up
    export $(grep -v '^#' $env_pkrvars_file | xargs)
}

export -f pad_string
export -f print_packer_variables
export -f get_build_environment