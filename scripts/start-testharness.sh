#!/bin/bash

# Verify that required env variables are set
if [ -z ${API_URI} ]; then
  echo "Environment variable API_URI missing."
  exit 1
fi

/testharness &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?