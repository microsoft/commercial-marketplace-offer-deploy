#!/bin/bash

version=$1

# parent mod
git tag -a $version -m "$version"  
git push origin $version

# sdk
sdk_version=sdk/$version
git tag -a $sdk_version -m "$sdk_version" 
git push origin $sdk_version

# sdk
tools_version=tools/$version
git tag -a $tools_version -m "$tools_version" 
git push origin $tools_version