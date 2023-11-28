import os
import json
from modm.marketplace.create_ui_definition import CreateUiDefinition


class CreateUiDefinitionInstallerStep:
    name = "installer"

    def __init__(self, step: dict, outputs: dict):
            """
            Initializes a new instance of the CreateUiDefinitionStep class.

            Args:
                document (dict): The UI definition step
            """
            self.step = step
            self.outputs = outputs

    def append_to(self, create_ui_definition: CreateUiDefinition):
        document = create_ui_definition.document
        steps = document["parameters"]["steps"]
        steps.append(self.step)

        outputs: dict = document["parameters"]["outputs"]
        outputs.update(self.outputs)


def from_file(file_path) -> CreateUiDefinitionInstallerStep:
    if not os.path.exists(file_path):
        raise FileNotFoundError(f"Could not find view definition file at {file_path}")
    
    with open(file_path, "r") as f:
        document = json.load(f)
        step = document["parameters"]["steps"][0]
        outputs = document["parameters"]["outputs"]
        return CreateUiDefinitionInstallerStep(step, outputs)