#!/bin/bash

package_file="./bin/app.zip"  # Change the path to the bin directory
main_template_file="./obj/mainTemplate.json"
create_ui_definition_file="./obj/createUiDefinition.json"
#terraform_content_dir="terraformContent"

# Change working directory to the vm directory
#cd "$(dirname "$0")"

# Create the zip file including the specified files and directories
zip -FS -j "$package_file" "$main_template_file" "$create_ui_definition_file" ./obj/content.zip
# find "$terraform_content_dir" -type f -print | zip -u "$package_file" -@

echo "Package app.zip created in the bin directory."
