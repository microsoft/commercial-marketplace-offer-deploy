import os
from pathlib import Path
import unittest

from modm.marketplace import MainTemplate, MainTemplateFinalizer
from modm.installer.installer_package_result import InstallerPackageResult
from modm.release.release_info import ReferenceInfo, OfferInfo
from tests import TestCaseBase


class TestMainTemplateFinalizer(TestCaseBase):
    def setUp(self):
        self.main_template = MainTemplate.from_file(self.data_path / "mainTemplate.json")

        self.reference_info: ReferenceInfo = ReferenceInfo.from_dict({
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
        })
        self.finalizer = MainTemplateFinalizer(self.main_template)

    def test_finalize_with_none_installer_resources(self):
        with self.assertRaises(AttributeError):
            self.finalizer.finalize()

    def test_finalize_with_false_use_vmi_reference(self):
        installer_package=InstallerPackageResult(self.data_path / "installer.zip")
        result = self.finalizer.finalize(reference_info=self.reference_info, use_vmi_reference=False, installer_package=installer_package)

        self.assertEqual(result.vm_offer, self.reference_info.offer.serialize())

