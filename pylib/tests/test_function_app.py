import logging
import unittest
from modm.marketplace.client_app import create_client_app_name

log = logging.getLogger(__name__)

class TestClientApp(unittest.TestCase):
    def test_create_client_app_name(self):
        prefix = "modm"
        client_app_name = create_client_app_name(prefix)
        log.info(client_app_name)
        self.assertTrue(client_app_name.startswith(prefix))
        self.assertEqual(len(client_app_name), len(prefix) + 12)