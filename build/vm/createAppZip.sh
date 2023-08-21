#!/bin/bash

package_file="app.zip"  # Replace with the actual file name

if [ -e "$package_file" ]; then
    echo "File exists. Deleting $package_file..."
    rm "$package_file"
    echo "File deleted."
else
    echo "File does not exist."
fi

zip $package_file mainTemplate.json createUiDefinition.json