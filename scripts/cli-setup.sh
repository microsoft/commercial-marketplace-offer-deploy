#!/bin/bash

echo "Setting up Development CLI..."
echo ""
cli_path=$PWD/cli

function modm() {
    cwd=$PWD
    pushd $cli_path &> /dev/null
        python -m devcli "$@" $cwd
    popd &> /dev/null
}

modm --help