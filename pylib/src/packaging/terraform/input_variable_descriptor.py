import re


class TerraformInputVariableDescriptor(object):
    def __init__(self, dict):
        self.name = list(dict.keys())[0]
        self.type = self.extract_type(list(dict.values())[0]['type'])

    def extract_type(self, type_string):
        pattern = r"\${(.+?)(?:\(.+\))?}"
        match = re.search(pattern, type_string)
        if match:
            return match.group(1)
        else:
            return type_string


