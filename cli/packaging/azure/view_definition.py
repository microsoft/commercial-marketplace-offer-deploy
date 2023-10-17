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
        rendered_json = self._template(self.inputs)
        return json.dumps(rendered_json, indent=4)
            

    @staticmethod
    def from_resource():
        resource_files = files("resources.templates")
        with as_file(resource_files.joinpath(ViewDefinition.file_name)) as resource_file:
            viewDefinition = ViewDefinition(resource_file)
            return viewDefinition
