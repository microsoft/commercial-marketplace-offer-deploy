import os
import unittest
from modm.arm.arm_template_parameter import ArmTemplateParameter
from modm.marketplace.create_ui_definition import CreateUiDefinition


class TestCreateUiDefinition(unittest.TestCase):
    def setUp(self):
        self.data_path = os.path.join(os.path.dirname(__file__), "data")
        self.file_path = os.path.join(self.data_path, "createUiDefinition.json")
        self.create_ui_definition = CreateUiDefinition.from_file(self.file_path)

    def test_validate(self):
        # Test case where outputs is None
        self.create_ui_definition.document = {"parameters": {"outputs": None}}

        result = self.create_ui_definition.validate(None)[0]
        self.assertTrue(isinstance(result, ValueError))
        self.assertEqual(str(result), "The createUiDefinition.json must contain an outputs section")

    def test_validate_outputs_and_input_parameters_not_matching(self):
        # Test case where outputs and template_input_parameters have different keys

        self.create_ui_definition.document = {"parameters": {"outputs": {"output1": "value1"}}}
        template_input_parameters = [
            ArmTemplateParameter(name="output1", type="string"),
            ArmTemplateParameter(name="output2", type="string"),
        ]

        result = self.create_ui_definition.validate(template_input_parameters)
        self.assertEqual(len(result), 1)
        self.assertTrue(isinstance(result[0], ValueError))

    def test_validate_outputs_match(self):
        # Test case where outputs and template_input_parameters have the same keys
        self.create_ui_definition.document = {"parameters": {"outputs": {"param1": "value1", "param2": "value2"}}}
        template_input_parameters = [
            ArmTemplateParameter(name="param1", type="string"),
            ArmTemplateParameter(name="param2", type="string"),
        ]

        result = self.create_ui_definition.validate(template_input_parameters)
        self.assertEqual((len(result)), 0)
    
    def test_validate_reserved_parameters(self):
        # Test case where outputs and template_input_parameters have the same keys
        self.create_ui_definition.document = {"parameters": {"outputs": {"resourceGroupName": "value1", "param1": "value2"}}}
        template_input_parameters = [
            ArmTemplateParameter(name="param1", type="string"),
        ]

        result = self.create_ui_definition.validate(template_input_parameters)
        self.assertEqual((len(result)), 1)
    
