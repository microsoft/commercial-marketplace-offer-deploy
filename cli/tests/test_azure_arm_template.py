import unittest
from modm.arm.arm_template import ArmTemplate


class TestArmTemplate(unittest.TestCase):
    def test_set_value_top_level_key(self):
        # Test setting a value for a top-level key
        template = ArmTemplate({"key1": "value1"})
        template.set_value("key2", "value2")

        expected = {"key1": "value1", "key2": "value2"}
        self.assertEqual(template.document, expected)

    def test_set_value_nested_key(self):
        # Test setting a value for a nested key
        template = ArmTemplate({"key1": {"key2": "value2"}})
        template.set_value("key1.key3", "value3")

        expected = {"key1": {"key2": "value2", "key3": "value3"}}
        self.assertEqual(template.document, expected)

    def test_set_value_non_existent_nested_key(self):
        # Test setting a value for a non-existent nested key
        template = ArmTemplate({"key1": {"key2": "value2"}})
        template.set_value("key1.key3.key4", "value4")

        expected = {"key1": {"key2": "value2", "key3": {"key4": "value4"}}}
        self.assertEqual(template.document, expected)

    def test_set_value_non_existent_top_level_key(self):
        # Test setting a value for a non-existent top-level key
        template = ArmTemplate({})
        template.set_value("key1.key2", "value2")

        expected = {"key1": {"key2": "value2"}}
        self.assertEqual(template.document, expected)

    def test_set_value_non_dict_parent(self):
        # Test setting a value for a key with a non-dict parent
        template = ArmTemplate({"key1": "value1"})
        with self.assertRaises(ValueError):
            template.set_value("key1.key2", "value2")
