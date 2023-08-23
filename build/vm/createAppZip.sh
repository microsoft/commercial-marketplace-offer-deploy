#!/bin/bash

package_file="app.zip"  # Replace with the actual file name
package_file_path=../../bin/$package_file

# Add the terraformcontent directory and its contents to the app.zip
zip -FS -r -j $package_file_path \
  "../../obj/mainTemplate.json" \
  "../../obj/createUiDefinition.json" \
  "terraformcontent/"
