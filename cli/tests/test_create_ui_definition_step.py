from copy import deepcopy
import os
import unittest
from packaging.installer.create_ui_definition_step import from_file
from packaging.azure.create_ui_definition import CreateUiDefinition

class TestCreateUiDefinitionInstallerStep(unittest.TestCase):
    def setUp(self):
        self.base_path = os.path.join(os.path.dirname(__file__))

        # pull the template file for the installer step located at the repo root
        self.template_file = os.path.join(self.base_path,  "../../templates", "createUiDefinition.json")
        self.file_path = os.path.join(self.base_path, "data", "createUiDefinition.json")

        self.create_ui_definition = CreateUiDefinition.from_file(self.file_path)

    def test_append_to(self):
        create_ui_definition = deepcopy(self.create_ui_definition)
        installer_step = from_file(self.template_file)

        installer_step.append_to(create_ui_definition)

        self.assertIn(installer_step.step, create_ui_definition.document["parameters"]["steps"])
        self.assertEqual(create_ui_definition.document["parameters"]["outputs"]["installer_installer"], "[steps('installer').username]")
        self.assertEqual(create_ui_definition.document["parameters"]["outputs"]["installer_password"], "[steps('installer').password]")