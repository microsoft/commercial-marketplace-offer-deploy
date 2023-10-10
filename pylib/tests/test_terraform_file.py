import unittest
import os
from packaging.terraform.terraform_file import TerraformFile

class TestTerraformFile(unittest.TestCase):
    def setUp(self):
        self.file_path = os.path.join(os.path.dirname(__file__), 'data/variables.tf')
        self.terraform_file = TerraformFile(self.file_path)

    def test_parse_variables(self):
        variables = self.terraform_file.parse_variables()
        self.assertEqual(len(variables), 6)
        self.assertEqual(variables[0].name, 'string_variable')
        self.assertEqual(variables[0].type, 'string')

        self.assertEqual(variables[1].name, 'bool_variable')
        self.assertEqual(variables[1].type, 'bool')
