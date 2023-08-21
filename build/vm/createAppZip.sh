#!/bin/bash

package_file="app.zip"  # Replace with the actual file name
package_file_path=../../bin/$package_file

zip -FS -r -j $package_file_path "../../obj/mainTemplate.json" "../../obj/createUiDefinition.json"