import unittest
from modm.installer.reserved_template_parameter import ReservedTemplateParameter, is_reserved

class TestReservedTemplateParameter(unittest.TestCase):
    def test_is_reserved(self):
        self.assertFalse(is_reserved("PARAMETER1"))
        self.assertTrue(is_reserved(ReservedTemplateParameter.resource_group_name))
        