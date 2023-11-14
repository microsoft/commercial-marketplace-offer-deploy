from importlib.resources import as_file, files
import os
import json
from pathlib import Path
from pybars import Compiler

from packaging.azure.create_ui_definition import CreateUiDefinition


class CreateUiDefinitionInstallerStep:
    name = "installer"
    output_value_format = "[steps('{}')]"

    def __init__(self, step: dict):
            """
            Initializes a new instance of the CreateUiDefinitionStep class.

            Args:
                document (dict): The UI definition step
            """
            self.step = step

    def append_to(self, create_ui_definition: CreateUiDefinition):
        document = create_ui_definition.document
        steps = document["parameters"]["steps"]
        steps.append(self.step)

        outputs = document["parameters"]["outputs"]

        # Add the output value for the installer step
        outputs["_installerUsername"] = "[steps('installer').username]"
        outputs["_installerPassword"] = "[steps('installer').password]"


def from_file(file_path) -> CreateUiDefinitionInstallerStep:
    if not os.path.exists(file_path):
        raise FileNotFoundError(f"Could not find view definition file at {file_path}")
    
    with open(file_path, "r") as f:
        document = json.load(f)
        step = document["parameters"]["steps"][0]
        return CreateUiDefinitionInstallerStep(step)