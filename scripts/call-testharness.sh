#!/bin/bash

testharness=http://localhost:8280
curl -s $testharness/
sleep 5
echo ""

echo "create event hook"
curl -s $testharness/createeventhook | jq .
sleep 5
echo ""

echo "create deployment"
curl -s $testharness/createdeployment | jq .
sleep 5
echo ""

echo "starting dry run"
curl -s $testharness/dryrun/1 | jq .
sleep 10
echo ""

echo "starting deployment"
curl -s $testharness/startdeployment/1 | jq .
sleep 20
echo ""

echo "retry for stage"
curl -s $testharness/redeploy/1/storageAccounts | jq .
sleep 10
echo ""
