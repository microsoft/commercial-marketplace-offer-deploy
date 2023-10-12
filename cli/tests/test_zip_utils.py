import os
import shutil
import tempfile
import unittest
from packaging.zip_utils import zip_dir, unzip_file

class testZipUtils(unittest.TestCase):
    def setUp(self):
        self.test_dir = tempfile.mkdtemp()

    def tearDown(self):
        shutil.rmtree(self.test_dir)

    def test_zip_operations(self):
        # test zip_dir
        dir_path = os.path.join(os.path.dirname(__file__), 'data')
        file_path = os.path.join(self.test_dir, 'test.pkg')

        zip_dir(dir_path, file_path)
        self.assertTrue(os.path.exists(file_path))

        unzip_file(file_path, self.test_dir)
        self.assertTrue(os.path.exists(os.path.join(self.test_dir, 'installer.pkg')))

        