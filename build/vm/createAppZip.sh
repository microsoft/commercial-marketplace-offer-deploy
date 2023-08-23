#!/bin/bash

package_file="app.zip"  # Replace with the actual file name
package_file_path="../../bin/$package_file"
terraform_content_dir="terraformContent"  # Name of the directory to include

# Create the zip file including the specified files and directory
zip -FS -j $package_file_path \
    "../../obj/mainTemplate.json" \
    "../../obj/createUiDefinition.json" \
    "$terraform_content_dir/*"   # Include the contents of the terraformcontent directory
