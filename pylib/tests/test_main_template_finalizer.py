import os
from pathlib import Path
import unittest
from modm.marketplace import MainTemplate, MainTemplateFinalizer
from modm.installer.installer_package_result import InstallerPackageResult
from modm.marketplace._resources import InstallerResources


class TestMainTemplateFinalizer(unittest.TestCase):
    def setUp(self):
        self.data_path = Path(os.path.join(os.path.dirname(__file__), "data"))
        self.main_template = MainTemplate.from_file(os.path.join(self.data_path, "mainTemplate.json"))

        release_reference = {"vmi": "test_id", "offer": { "plan": {}, "imageReference": {}}}
        self.installer_resources = InstallerResources(self.data_path, "latest", release_reference)

        self.finalizer = MainTemplateFinalizer(self.main_template)

    def test_finalize_with_none_installer_resources(self):
        with self.assertRaises(ValueError):
            self.finalizer.finalize()

    def test_finalize_with_false_use_vmi_reference(self):
        installer_package=InstallerPackageResult(self.data_path.joinpath("installer.zip"))
        result = self.finalizer.finalize(installer_resources=self.installer_resources, installer_package=installer_package)
        self.assertEqual(result.vm_offer, self.installer_resources.vm_offer)

