#! /bin/bash

output=$(docker rmi modm 2>&1)   # execute the command and capture output

if [[ $output =~ container\ ([a-f0-9]{12})\ is\ using\ its\ referenced\ image ]]; then
    container_id=${BASH_REMATCH[1]}
    echo "Container ID: $container_id"
    output=$(docker rm $container_id 2>&1) 
    echo $output
    output=$(docker rmi modm 2>&1)
    echo $output
else
    echo "Container ID not found in output."
    echo $output
fi

output=$(docker rmi testharness 2>&1)   # execute the command and capture output

if [[ $output =~ container\ ([a-f0-9]{12})\ is\ using\ its\ referenced\ image ]]; then
    container_id=${BASH_REMATCH[1]}
    echo "Container ID: $container_id"
    output=$(docker rm $container_id 2>&1) 
    echo $output
    output=$(docker rmi testharness 2>&1)
    echo $output
else
    echo "Container ID not found in output."
    echo $output
fi

echo "Removing tmp directory"
output=$(rm -rf ~/tmp 2>&1) 
echo $output

