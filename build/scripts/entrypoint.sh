#!/bin/bash


# Start the api server
/apiserver & >&1

# Start the operator server
/operator & >&1 

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?