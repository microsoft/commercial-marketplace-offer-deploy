#!/bin/bash

# Python SDK build

# docs: https://pypa-build.readthedocs.io/en/latest/
# copy files to resources folder to include them in the distributable
echo "Building Python SDK."

# prefix output with indent
exec > >(trap "" INT TERM; sed 's/^/  /')
exec 2> >(trap "" INT TERM; sed 's/^/   (stderr) /' >&2)

echo "copying files to resources folder."

cp -r ./templates ./sdks/python/resources

out_dir=$PWD/dist

pushd ./sdks/python
  echo "Installing build tools."
  python -m pip install --upgrade build

  echo "Building wheel..."
  python -m build --wheel --outdir $out_dir
  echo ""
popd