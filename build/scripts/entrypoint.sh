#!/bin/bash

echo "MODM starting in 3 minutes..."
sleep 180

echo "API Server starting."
# Start the api server
/apiserver /dev/fd/1 2>&1 &

echo "Operator starting."
# Start the operator server
/operator /dev/fd/1 2>&1 &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?