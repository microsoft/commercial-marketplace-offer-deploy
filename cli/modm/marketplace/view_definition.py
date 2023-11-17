from importlib.resources import as_file, files
import os
import json
from pathlib import Path
from pybars import Compiler


class ViewDefinition:
    file_name = "viewDefinition.json"

    def __init__(self, file: Path):
        self._compiler = Compiler()
        self.file = file
        self.name = self.file_name
        self.inputs = {}

    def add_input(self, name, value):
        self.inputs[name] = value

    def to_json(self):
        # read file to string
        with open(self.file, "r") as f:
            self._template = f.read()
        
        self._template = self._compiler.compile(self._template)
        rendered_json = json.loads(self._template(self.inputs))
        return json.dumps(rendered_json, indent=4)


def from_file(file_path) -> ViewDefinition:
    if not os.path.exists(file_path):
        raise FileNotFoundError(f"Could not find view definition file at {file_path}")
    return ViewDefinition(file_path)