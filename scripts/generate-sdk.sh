#! /bin/bash

# Call this scripts while in the ./api directory
cd ./api
autorest autorest --go --go-sdk-folder=../sdk/internal/generated --tag=preview-2023-03-01