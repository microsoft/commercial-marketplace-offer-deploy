# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
# from pathlib import Path
# import shutil
# import tempfile
# import unittest
# from unittest.mock import MagicMock, patch
# from modm.marketplace.application_packaging_options import ApplicationPackageOptions
# from modm.release.resources_archive import ResourcesArchive
# from modm.release.version import Version


# class TestApplicationPackageOptions(unittest.TestCase):
#     def test_init(self):
#         options = ApplicationPackageOptions("1.0.0", True, "vmi-id", "resources.json", "/tmp")
#         self.assertEqual(options.version, "1.0.0")
#         self.assertEqual(options.use_vmi_reference, True)
#         self.assertEqual(options.vmi_reference_id, "vmi-id")
#         self.assertEqual(options.resources_archive, Path("resources.json"))
#         self.assertEqual(options.out_dir, "/tmp")

#     def test_out_dir(self):
#         options = ApplicationPackageOptions("1.0.0")
#         self.assertIsNotNone(options.out_dir)

#     def test_resources_file(self):
#         options = ApplicationPackageOptions("1.0.0", resources_archive_file="resources.json")
#         self.assertEqual(options.resources_archive, Path("resources.json"))

#     def test_installer_version_latest(self):
#         options = ApplicationPackageOptions("latest")
#         self.assertIsNotNone(options.version)

#     def test_installer_version_str(self):
#         options = ApplicationPackageOptions("1.0.0")
#         self.assertEqual(options.version, "1.0.0")

#     def test_installer_version_version(self):
#         options = ApplicationPackageOptions(Version("1.0.0"))
#         self.assertEqual(options.version, Version("1.0.0"))

#     @patch("modm.release.release_provider.ReleaseProvider.get")
#     @patch("modm.release.release_provider.ReleaseProvider.get_resources")
#     def test_vmi_reference_id(self, mock_get_resources, mock_get):
#         mock_get_resources.return_value = MagicMock()
#         mock_get.return_value = MagicMock()

#         options = ApplicationPackageOptions("1.0.0", vmi_reference_id="vmi-id")
#         self.assertEqual(options.vmi_reference_id, "vmi-id")

#     @patch("modm.release.release_provider.ReleaseProvider.get")
#     @patch("modm.release.release_provider.ReleaseProvider.get_resources")
#     def test_use_vmi_reference(self, mock_get_resources, mock_get):
#         mock_get_resources.return_value = MagicMock()
#         mock_get.return_value = MagicMock()

#         options = ApplicationPackageOptions("1.0.0", use_vmi_reference=True)
#         self.assertEqual(options.use_vmi_reference, True)