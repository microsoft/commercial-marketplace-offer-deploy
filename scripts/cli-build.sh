#!/bin/bash

# Python SDK build

# docs: https://pypa-build.readthedocs.io/en/latest/
# copy files to resources folder to include them in the distributable
echo "Building Python Library."

echo "copying files to resources folder."

echo "Installing build tools."
python -m pip install --upgrade build

echo "Executing tests"
python -m unittest discover ./cli -v

echo "Building wheel..."
python -m build --wheel ./cli