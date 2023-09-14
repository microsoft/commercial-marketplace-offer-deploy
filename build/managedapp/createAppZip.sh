#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Usage: $0 <scenario_name>"
  exit 1
fi

scenario_name="$1"
package_file="./bin/app.zip"  # Change the path to the bin directory
main_template_file="./obj/mainTemplate.json"
create_ui_definition_file="./obj/createUiDefinition.json"

# Create the content.zip file
zip -r ./obj/content.zip ./build/managedapp/$scenario_name/*

# Create the zip file including the specified files and directories
zip -FS -j "$package_file" "$main_template_file" "$create_ui_definition_file" ./obj/content.zip

echo "Package app.zip created in the bin directory."
