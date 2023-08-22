#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Usage: $0 <artifacts location>"
  exit 1
fi

ARTIFACTS_LOCATION="$1"

echo "Hello, world from script extension!  The _artifactsLocation is $ARTIFACTS_LOCATION" > scriptExtensionOutput.txt