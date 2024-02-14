# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
from copy import deepcopy
import os
import unittest
from modm.marketplace.create_ui_definition_step import from_file
from modm.marketplace.create_ui_definition import CreateUiDefinition
from tests import TestCaseBase

class TestCreateUiDefinitionInstallerStep(TestCaseBase):
    def setUp(self):
        self.base_path = os.path.join(os.path.dirname(__file__))

        # pull the template file for the installer step located at the repo root
        self.template_file = os.path.join(self.base_path,  "../../templates", "createUiDefinition.json")
        self.file_path = self.data_path / "createUiDefinition.json"

        self.create_ui_definition = CreateUiDefinition.from_file(self.file_path)

    def test_append_to(self):
        create_ui_definition = deepcopy(self.create_ui_definition)
        installer_step = from_file(self.template_file)

        installer_step.append_to(create_ui_definition)

        self.assertIn(installer_step.step, create_ui_definition.document["parameters"]["steps"])

        outputs = create_ui_definition.document["parameters"]["outputs"]

        for key, _ in installer_step.outputs.items():
            self.assertEqual(outputs[key], installer_step.outputs[key])
    
    def test_append_to_should_not_overwrite_existing_output_value(self):
        create_ui_definition = deepcopy(self.create_ui_definition)
        installer_step = from_file(self.template_file)
        installer_step.outputs["location"] = "overridden value"

        installer_step.append_to(create_ui_definition)

        self.assertIn("location", create_ui_definition.document["parameters"]["outputs"])
        self.assertNotEqual(create_ui_definition.document["parameters"]["outputs"]["location"], "overridden value")