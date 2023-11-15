from pathlib import Path
import unittest
from modm.marketplace.application_packaging_options import ApplicationPackageOptions
from modm.release.version import Version


class TestApplicationPackageOptions(unittest.TestCase):
    def test_init(self):
        options = ApplicationPackageOptions("1.0.0", True, "vmi-id", "resources.json", "/tmp")
        self.assertEqual(options.version, "1.0.0")
        self.assertEqual(options.use_vmi_reference, True)
        self.assertEqual(options.vmi_reference_id, "vmi-id")
        self.assertEqual(options.resources_file, Path("resources.json"))
        self.assertEqual(options.out_dir, "/tmp")

    def test_out_dir(self):
        options = ApplicationPackageOptions("1.0.0")
        self.assertIsNotNone(options.out_dir)

    def test_resources_file(self):
        options = ApplicationPackageOptions("1.0.0", resources_file="resources.json")
        self.assertEqual(options.resources_file, Path("resources.json"))

    def test_installer_version_latest(self):
        options = ApplicationPackageOptions("latest")
        self.assertIsNotNone(options.version)

    def test_installer_version_str(self):
        options = ApplicationPackageOptions("1.0.0")
        self.assertEqual(options.version, "1.0.0")

    def test_installer_version_version(self):
        options = ApplicationPackageOptions(Version("1.0.0"))
        self.assertEqual(options.version, Version("1.0.0"))

    def test_vmi_reference_id(self):
        options = ApplicationPackageOptions("1.0.0", vmi_reference_id="vmi-id")
        self.assertEqual(options.vmi_reference_id, "vmi-id")

    def test_use_vmi_reference(self):
        options = ApplicationPackageOptions("1.0.0", vmi_reference=True)
        self.assertEqual(options.use_vmi_reference, True)