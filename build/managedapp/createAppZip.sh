#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Usage: $0 <scenario_name>"
  exit 1
fi

scenario_name="$1"
echo "The scenario name is: $scenario_name"
package_file="./bin/app.zip"  # Change the path to the bin directory
main_template_file="./obj/mainTemplate.json"
create_ui_definition_file="./obj/createUiDefinition.json"

origional_dir=$(pwd)
# Change the directory to the scenario's location
echo "Changing directory to ./build/managedapp/$scenario_name"
cd "./build/managedapp/$scenario_name" || exit
echo "The current directory is: $(pwd)"

# Create the content.zip file
zip -r "../../../../obj/content.zip" ./*
echo "zipped content.zip"

# Go back to the original directory
echo "going back to origional directory"
cd -
echo "The current directory is: $(pwd)"
echo "The ./obj directory contains: $(ls -la ./obj)"

# Create the zip file including the specified files and directories
zip -FS -j "$package_file" "$main_template_file" "$create_ui_definition_file" ./obj/content.zip

echo "Package app.zip created in the bin directory."
