#!/bin/bash

version=$1

echo "versioning: $version"

# sdk
cd ./sdk
sdk_version=sdk/$version
go mod tidy

git add ./go.mod ./go.sum
git commit -m "$sdk_version"  
git push origin v1 -f

git tag -a $sdk_version -m "$sdk_version" 
git push origin $sdk_version


cd ../
go get github.com/microsoft/commercial-marketplace-offer-deploy/sdk@none
go get github.com/microsoft/commercial-marketplace-offer-deploy/sdk@$version

go mod tidy
git add ./go.mod ./go.sum
git commit -m "updating tools to consume $version"  
git push origin main -f

# git tag -a $tools_version -m "$tools_version" 
# git push origin $tools_version