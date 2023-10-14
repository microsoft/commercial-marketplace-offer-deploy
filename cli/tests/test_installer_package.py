import shutil
import tempfile
import unittest
import os
from packaging import ManifestInfo, InstallerPackage, DeploymentType

class TestInstallerPackage(unittest.TestCase):
    def setUp(self):
        self.data_path = os.path.join(os.path.dirname(__file__), 'data')
        self.main_template_file = os.path.join(self.data_path, 'simple_terraform', 'main.tf')

        self.manifest = ManifestInfo(main_template=self.main_template_file, 
                                     deployment_type=DeploymentType.terraform)

        self.installer_package = InstallerPackage(self.manifest)

    def test_file_name(self):
        self.assertEqual(InstallerPackage.file_name, 'installer.pkg')

    def test_init(self):
        self.assertEqual(self.installer_package.manifest.main_template, self.manifest.main_template)

    def test_create(self):
        file = self.installer_package.create()
        self.assertTrue(file.exists())

        # now unpack and verify
        temp_dir = tempfile.mkdtemp()
        self.installer_package.unpack(file, temp_dir)

        self.assertTrue(os.path.exists(os.path.join(temp_dir, 'main.tf')))

        # clean up
        shutil.rmtree(file.parent)
        shutil.rmtree(temp_dir)