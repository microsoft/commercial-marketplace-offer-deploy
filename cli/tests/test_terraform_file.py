import unittest
import os
from modm.terraform.terraform_file import TerraformFile
from tests import TestCaseBase

class TestTerraformFile(TestCaseBase):
    def setUp(self):
        self.file_path = self.data_path / 'variables.tf'
        self.terraform_file = TerraformFile(self.file_path)

    def test_parse_variables(self):
        variables = self.terraform_file.parse_variables()
        self.assertEqual(len(variables), 8)

        self.assertEqual(variables[0].name, 'string_variable')
        self.assertEqual(variables[0].type, 'string')

        self.assertEqual(variables[1].name, 'string_variable_sensitive')
        self.assertEqual(variables[1].type, 'string')

        self.assertEqual(variables[2].name, 'bool_variable')
        self.assertEqual(variables[2].type, 'bool')

        self.assertEqual(variables[3].name, 'number_variable')
        self.assertEqual(variables[3].type, 'number')

        self.assertEqual(variables[4].name, 'list_variable')
        self.assertEqual(variables[4].type, 'list')
        
        self.assertEqual(variables[5].name, 'set_variable')
        self.assertEqual(variables[5].type, 'set')

        self.assertEqual(variables[6].name, 'object_variable')
        self.assertEqual(variables[6].type, 'object')

        self.assertEqual(variables[7].name, 'map_variable')
        self.assertEqual(variables[7].type, 'map')


    def test_reserved_input_parameter(self):
        file_path = self.data_path / 'reserved_parameters/main.tf'
        file = TerraformFile(file_path)

        # resource group name should exists in the result
        variables = file.parse_variables()
        self.assertEqual(len(variables), 2)