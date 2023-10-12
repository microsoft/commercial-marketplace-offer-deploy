import re
from packaging.terraform.variable_types import VariableType

class TerraformInputVariableDescriptor(object):
    def __init__(self, dict):
        self.name = list(dict.keys())[0]
        self.type = self.extract_type(list(dict.values())[0]['type'])

        if not VariableType.is_valid_type(self.type):
            raise ValueError("Invalid type: {}".format(self.type))

    def extract_type(self, type_string):
        pattern = r"\${(.+?)(?:\(.+\))?}"
        match = re.search(pattern, type_string)
        if match:
            type = match.group(1)
            return type.replace('()', '')
        else:
            return type_string
