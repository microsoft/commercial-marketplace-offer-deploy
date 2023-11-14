import unittest
import os
from modm.azure.arm_template import ArmTemplate, ArmTemplateParameter

class TestArmTemplate(unittest.TestCase):
    def setUp(self):
        self.file_path = os.path.join(os.path.dirname(__file__), 'data/mainTemplate.json')
        self.arm_template = ArmTemplate.from_file(self.file_path)

    def test_set_parameter(self):
        parameter = ArmTemplateParameter("param1", "string")
        self.arm_template.set_parameter(parameter)

        parameters = self.arm_template.document["parameters"]
        self.assertEqual(parameters["param1"]["type"], "string")