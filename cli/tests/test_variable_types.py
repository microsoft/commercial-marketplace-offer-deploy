import unittest
from modm.terraform.variable_types import TerraformInputVariableType

class TestVariableType(unittest.TestCase):
    def test_to_list(self):
        expected = ['string', 'bool', 'number', 'list', 'set', 'object', 'map']
        self.assertEqual(TerraformInputVariableType.to_list(), expected)

    def test_is_valid_type(self):
        expected = ['string', 'bool', 'number', 'list', 'set', 'object', 'map']
        for type in expected:
            self.assertTrue(TerraformInputVariableType.is_valid_type(type))
        
        self.assertFalse(TerraformInputVariableType.is_valid_type("int"))