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

original_dir=$(pwd)
echo "The original directory is: $original_dir"
# Change the directory to the scenario's location
echo "Changing directory to ./build/managedapp/$scenario_name"
cd "./build/managedapp/$scenario_name" || exit
echo "The current directory is: $(pwd)"

echo "current value of the .obj directory is $(ls -la $original_dir/obj)"

# Create the content.zip file
zip -r "$original_dir/obj/content.zip" ./content/*
echo "zipped content.zip"

# Go back to the original directory
echo "going back to original directory"
cd -
echo "The current directory is: $(pwd)"
echo "The ./obj directory contains: $(ls -la ./obj)"

# Create the hash value for the content.zip file
hash_value=$(openssl dgst -sha256 "./obj/content.zip" | awk '{print $2}')
echo "SHA256 hash for ./obj/content.zip: $hash_value"

TEMP_FILE="./obj/mainTemplateUpdate.json"
echo "removing TEMP_FILE: $TEMP_FILE"
rm $TEMP_FILE 2> /dev/null
echo "$(ls -la ./obj)"
echo "creating TEMP_FILE with replacement: $TEMP_FILE"
sed "s|<CONTENT_SIG>|$hash_value|g" "$main_template_file" > "$TEMP_FILE"
echo "$(ls -la ./obj)"

echo "removing main_template_file: $main_template_file"
rm "$main_template_file" 2> /dev/null
echo "$(ls -la ./obj)"

echo "copying TEMP_FILE to main_template_file: $main_template_file"
cp -f "$TEMP_FILE" "$main_template_file"
echo "$(ls -la ./obj)"

# Create the zip file including the specified files and directories
zip -FS -j "$package_file" "$main_template_file" "$create_ui_definition_file" ./obj/content.zip

echo "Package app.zip created in the bin directory."
