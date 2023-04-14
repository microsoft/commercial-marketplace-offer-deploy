#!/bin/bash

test ${AZURE_STORAGE_MOUNT_POINT}
#rm -rf ${AZURE_STORAGE_MOUNT_POINT}
#mkdir ${AZURE_STORAGE_MOUNT_POINT}

##blobfuse ${AZURE_STORAGE_MOUNT_POINT} --use-https=true --tmp-path=/tmp/blobfuse/${AZURE_STORAGE_ACCOUNT} --container-name=${AZURE_STORAGE_ACCOUNT_CONTAINER} -o allow_other

# Start the api server
./apiserver & 

# Start the operator server
./operator &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?