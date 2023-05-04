#!/bin/bash

go build -o ./bin/ ./tools/testharness

mkdir -p ./bin/templates/
cp -n ./test/testdata/taggeddeployment/mainTemplateBicep.json ./bin/templates/mainTemplateBicep.json
cp -n ./test/testdata/taggeddeployment/parametersBicep.json ./bin/templates/parametersBicep.json

echo "Test Harness built"