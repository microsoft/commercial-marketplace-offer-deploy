#!/bin/bash

case_name=$1

if [ "$case_name" == "" ]; then
    case_name=success
fi

testharness=http://localhost:8280
curl -s $testharness/
sleep 5
echo ""

echo "setting case"
curl -s $testharness/setcase/$case_name
sleep 5
echo ""

echo "create event hook"
curl -s $testharness/createeventhook | jq .
sleep 5
echo ""

echo "create deployment"
resp=$(curl -s $testharness/createdeployment | jq . -r)
echo $resp
deployment_id=$(echo $resp | jq .id -r)
sleep 5
echo ""

echo "starting dry run [$deployment_id]"
curl -s $testharness/dryrun/$deployment_id | jq .
sleep 10
echo ""

echo "starting deployment [$deployment_id]"
curl -s $testharness/startdeployment/$deployment_id | jq .
sleep 20
echo ""

echo "retry for stage"
curl -s $testharness/redeploy/$deployment_id/storageAccounts | jq .
sleep 10
echo ""
