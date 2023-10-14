#!/bin/bash

# Python SDK build

# docs: https://pypa-build.readthedocs.io/en/latest/
# copy files to resources folder to include them in the distributable
echo "Building Python CLI."

# prefix output with indent
exec > >(trap "" INT TERM; sed 's/^/  /')
exec 2> >(trap "" INT TERM; sed 's/^/   (stderr) /' >&2)

echo "copying files to resources folder."

mkdir -p ./cli/resources
touch ./cli/resources/__init__.py
cp -r ./templates ./cli/resources
cp -r ./schemas ./cli/resources

echo "Installing build tools."
python -m pip install --upgrade build

echo "Executing tests"
python -m unittest discover ./cli -v

echo "Building wheel..."
python -m build --wheel ./cli