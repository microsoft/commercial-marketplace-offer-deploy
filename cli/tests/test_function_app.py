import logging
import unittest
from packaging.azure.function_app import create_client_app_name

log = logging.getLogger(__name__)

class TestFunctionApp(unittest.TestCase):
    def test_create_client_app_name(self):
        prefix = "modmfunc"
        client_app_name = create_client_app_name(prefix)
        log.info(client_app_name)
        self.assertTrue(client_app_name.startswith(prefix))
        self.assertEqual(len(client_app_name), len(prefix) + 12)