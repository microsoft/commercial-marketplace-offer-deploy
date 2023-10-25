import re
from .variable_types import TerraformInputVariableType

class TerraformInputVariable(object):
    def __init__(self, dict):
        attrs = list(dict.values())[0]

        self.name = list(dict.keys())[0]

        if 'type' not in attrs:
            raise ValueError("Missing type for variable: {} in the Terraform template.".format(self.name))
        
        self.type = self._extract_type(attrs['type'])
        self.sensitive = self._get_sensitive(attrs)
        
        if not TerraformInputVariableType.is_valid_type(self.type):
            raise ValueError("Invalid type: {}".format(self.type))

    def _get_sensitive(self, attrs):
        sensitive = False

        if attrs is None:
            return sensitive
        if 'sensitive' in attrs:
            sensitive = bool(attrs['sensitive'])
        return sensitive

    def _extract_type(self, type_string):
        pattern = r"\${(.+?)(?:\(.+\))?}"
        match = re.search(pattern, type_string)
        if match:
            type = match.group(1)
            return type.replace('()', '')
        else:
            return type_string
