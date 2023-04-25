#!/bin/bash


# Start the api server
/apiserver & 1> /proc/1/fd/1 2> /proc/1/fd/2 

# Start the operator server
/operator & 1> /proc/1/fd/1 2> /proc/1/fd/2 

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?