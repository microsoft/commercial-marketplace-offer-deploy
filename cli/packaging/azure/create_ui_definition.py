import json
import os

from .arm_template_parameter import ArmTemplateParameter


class CreateUiDefinition:
    def __init__(self, document):
        self.document = document

    def validate(self, template_parameters: list[ArmTemplateParameter]):
        validation_results = []
        outputs = self.document["parameters"]["outputs"]

        if outputs is None:
            validation_results.append(ValueError("The createUiDefinition.json must contain an outputs section"))
            return validation_results

        outputs_keys = set(outputs.keys())
        template_parameters_keys = set(list(map(lambda p: p.name, template_parameters)))
        diff = template_parameters_keys.symmetric_difference(outputs_keys)

        if len(diff) > 0:
            validation_results.append(
                ValueError(
                    { 
                    "message": "The outputs defined in createUiDefinition.json do not match the input parameters of your template.", 
                    "properties": list(diff) 
                    })
            )

        return validation_results

    def to_json(self):
        return json.dumps(self.document, indent=4)

    @staticmethod
    def from_file(file_path):
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Could not find create ui definition file at {file_path}")

        with open(file_path, "r") as f:
            document = json.load(f)
            return CreateUiDefinition(document)
