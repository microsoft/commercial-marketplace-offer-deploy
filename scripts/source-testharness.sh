#!/bin/bash

export testharness_url=http://localhost:8280

function fresh() {
    echo "Cleaning out"
    rm -rf ./*.txt
    rm -rf ./*.db
}


function testharness() {
    curl -s $testharness_url/
    echo ""
}

function setcase() {
    case_name=$1
    echo "Setting case: $case_name"
    curl -s $testharness_url/setcase/$case_name
    echo ""
}

function eventhook() {
    echo "Creating event hook"
    curl -s $testharness_url/createeventhook | jq .
    echo ""
}

function create() {
    echo "Creating deployment"
    resp=$(curl -s $testharness_url/createdeployment | jq . -r)
    echo $resp | jq .
    export deployment_id=$(echo $resp | jq .id -r)
    echo ""
}

function dryrun() {
    id=$1
    if [[ -z "$id" ]]; then
        id=$deployment_id
    fi
    echo "starting dry run [$id]"
    curl -s $testharness_url/dryrun/$id | jq .
    echo ""
}

function deploy() {
    id=$1
    if [[ -z "$id" ]]; then
        id=$deployment_id
    fi
    echo "starting deployment [$id]"
    curl -s $testharness_url/startdeployment/$id | jq .
    echo ""
}