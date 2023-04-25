#!/bin/bash


# Start the api server
/apiserver & 

# Start the operator server
/operator  & 

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?