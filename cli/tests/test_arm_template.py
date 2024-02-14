# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
import unittest
import os
from modm.arm.arm_template import ArmTemplate, ArmTemplateParameter
from tests import TestCaseBase

class TestArmTemplate(TestCaseBase):
    def setUp(self):
        self.file_path = self.data_path / 'mainTemplate.json'
        self.arm_template = ArmTemplate.from_file(self.file_path)

    def test_set_parameter(self):
        parameter = ArmTemplateParameter("param1", "string")
        self.arm_template.set_parameter(parameter)

        parameters = self.arm_template.document["parameters"]
        self.assertEqual(parameters["param1"]["type"], "string")