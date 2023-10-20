import logging
import unittest
from packaging.azure.function_app import create_function_app_name

log = logging.getLogger(__name__)

class TestFunctionApp(unittest.TestCase):
    def test_create_function_app_name(self):
        prefix = "modmfunc"
        function_app_name = create_function_app_name(prefix)
        log.info(function_app_name)
        self.assertTrue(function_app_name.startswith(prefix))
        self.assertEqual(len(function_app_name), len(prefix) + 12)