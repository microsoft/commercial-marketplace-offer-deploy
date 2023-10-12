import unittest
from packaging.terraform.variable_types import VariableType

class TestVariableType(unittest.TestCase):
    def test_to_list(self):
        expected = ['string', 'bool', 'number', 'list', 'set', 'object', 'map']
        self.assertEqual(VariableType.to_list(), expected)

    def test_is_valid_type(self):
        expected = ['string', 'bool', 'number', 'list', 'set', 'object', 'map']
        for type in expected:
            self.assertTrue(VariableType.is_valid_type(type))
        
        self.assertFalse(VariableType.is_valid_type("int"))