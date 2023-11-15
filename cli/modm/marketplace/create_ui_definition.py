from copy import deepcopy
import json
import os
from modm.installer.reserved_template_parameter import ReservedTemplateParameter
from ..arm.arm_template_parameter import ArmTemplateParameter


class CreateUiDefinition:
    def __init__(self, document):
        self.document = document

    def validate(self, template_parameters: list[ArmTemplateParameter]):
        reserved_template_parameters = ReservedTemplateParameter.all()
        validation_results = []
        outputs = deepcopy(self.document["parameters"]["outputs"])

        if outputs is None:
            validation_results.append(ValueError("The createUiDefinition.json must contain an outputs section"))
            return validation_results
        
        for reserved_param in reserved_template_parameters:
            if reserved_param in outputs:
                del outputs[reserved_param]
                validation_results.append(
                    ValueError(
                        {
                            "message": f"The outputs defined in createUiDefinition.json contain a reserved parameter: {reserved_param}",
                            "properties": [reserved_param]
                        }
                    )
                )

        outputs_keys = set(outputs.keys())
        template_parameters_keys = set(list(map(lambda p: p.name, template_parameters)))
        
        for reserved_key in reserved_template_parameters:
            if reserved_key in template_parameters_keys:
                template_parameters_keys.remove(reserved_key)

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
