#!/bin/bash

version=$1

echo "versioning: $version"

# sdk
cd ./sdk
sdk_version=sdk/$version

git tag -a $sdk_version -m "$sdk_version" 
git push origin $sdk_version


cd ../

# publish the root/parent
go get github.com/microsoft/commercial-marketplace-offer-deploy/sdk@$version
go mod tidy

git .
git commit -m "$version"  
git push origin main -f

git tag -a $version
git push origin $version