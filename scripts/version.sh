#!/bin/bash

version=$1

# parent mod
go mod tidy
git commit -m "$version"  
git tag -a $version -m "$version"  
git push origin $version

# sdk
cd ./sdk
sdk_version=sdk/$version
go get github.com/microsoft/commercial-marketplace-offer-deploy@none
go get github.com/microsoft/commercial-marketplace-offer-deploy@$version
go mod tidy

git add ./go.mod ./go.sum
git commit -m "$sdk_version"  
git push origin main -f

git tag -a $sdk_version -m "$sdk_version" 
git push origin $sdk_version

# tools
cd ../tools
go mod tidy

tools_version=tools/$version
go get github.com/microsoft/commercial-marketplace-offer-deploy@none
go get github.com/microsoft/commercial-marketplace-offer-deploy@$version
go get github.com/microsoft/commercial-marketplace-offer-deploy/sdk@$version

go mod tidy
git add ./go.mod ./go.sum
git commit -m "$tools_version"  
git push origin main -f

git tag -a $tools_version -m "$tools_version" 
git push origin $tools_version