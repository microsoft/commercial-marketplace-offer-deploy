# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
import os
from pathlib import Path
import unittest

from modm.marketplace import MainTemplate, MainTemplateFinalizer
from modm.installer.installer_package_result import InstallerPackageResult
from modm.release.release_info import ReleaseInfo
from tests import TestCaseBase


class TestMainTemplateFinalizer(TestCaseBase):
    def setUp(self):
        self.main_template = MainTemplate.from_file(self.data_path / "mainTemplate.json")

        self.release_info = ReleaseInfo.from_dict({
            "version": "1.0.0",
            "reference": {
                "vmi": "test_id", 
                "offer": { 
                    "plan": { "name": "", "publisher": "", "product": ""}, 
                    "imageReference": {
                        "publisher": "test_publisher",
                        "offer": "test_offer",
                        "sku": "test_sku",
                        "version": "test_version"
                    }
                }
            }
        })
        self.finalizer = MainTemplateFinalizer(self.main_template)

    def test_finalize_with_none_installer_resources(self):
        with self.assertRaises(AttributeError):
            self.finalizer.finalize()

    def test_finalize_with_false_use_vmi_reference(self):
        installer_package=InstallerPackageResult(self.data_path / "installer.zip")
        result = self.finalizer.finalize(release_info=self.release_info, use_vmi_reference=False, installer_package=installer_package)

        self.assertEqual(result.vm_offer, self.release_info.reference.offer.serialize())

