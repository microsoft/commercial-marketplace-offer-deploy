#!/bin/bash

go build -o ./bin/ ./tools/testharness

mkdir -p ./bin/template/
cp -n ./test/testdata/taggeddeployment/mainTemplateBicep.json ./bin/template/mainTemplateBicep.json
cp -n ./test/testdata/taggeddeployment/parametersBicep.json ./bin/template/parametersBicep.json

echo "Test Harness built"