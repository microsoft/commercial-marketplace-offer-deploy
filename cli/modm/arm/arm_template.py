import os
import json

from .arm_template_parameter import ArmTemplateParameter


class ArmTemplate:
    default_name = "mainTemplate.json"

    def __init__(self, document, name=None):
        self.document = document
        self.name = name or self.default_name

    def write(self, file_path):
        # Create directory if it doesn't exist
        os.makedirs(os.path.dirname(file_path), exist_ok=True)

        # Write formatted JSON to file
        with open(file_path, "w") as f:
            json.dump(self.document, f, indent=4)

    def set_value(self, key, value):
        self._set_value(self.document, key, value)
        return value

    def set_parameters(self, parameters: list[ArmTemplateParameter]):
        for parameter in parameters:
            self.set_parameter(parameter)

    def set_parameter(self, parameter: ArmTemplateParameter):
        self.document["parameters"][parameter.name] = parameter.value()

    def get_parameters(self) -> list[ArmTemplateParameter]:
        parameters = []
        for name, parameter in self.document["parameters"].items():
            parameters.append(ArmTemplateParameter(name, parameter['type']))

        return parameters

    def to_json(self):
        return json.dumps(self.document, indent=4)

    @staticmethod
    def from_file(file_path):
        """
        Load an ARM template from a file.

        Args:
          file_path (str): The path to the ARM template file.

        Returns:
          ArmTemplate: An instance of the ArmTemplate class representing the loaded ARM template.
        """
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Could not find ARM template file at {file_path}")

        with open(file_path, "r") as f:
            document = json.load(f)
            return ArmTemplate(document)
    
    def _set_value(self, document, key, value):
        if key is None:
            raise ValueError("key cannot be None")

        if "." in key:
            keys = key.split(".")
            current_key = keys.pop(0)
            current_value = document.get(current_key, None)

            if current_value is None:
                current_value = {}
                document[current_key] = current_value

            if isinstance(current_value, dict):
                if len(keys) > 0:
                    document = document[current_key]
                next_key = ".".join(keys)
                self._set_value(document, next_key, value)
            else:
                raise ValueError(f"Could not set value for key {current_key}")
        else:
            document[key] = value
