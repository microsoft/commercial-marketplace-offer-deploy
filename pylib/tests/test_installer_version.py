import unittest
from modm.installer import InstallerVersion

class TestInstallerVersion(unittest.TestCase):
    def test_constructor_with_version_name(self):
        # Test constructor with version name
        version = InstallerVersion("v1.2.3")
        self.assertEqual(version.major, 1)
        self.assertEqual(version.minor, 2)
        self.assertEqual(version.patch, 3)
        self.assertEqual(version.suffix, "")

    def test_constructor_with_major_minor_patch(self):
        # Test constructor with major, minor, patch
        version = InstallerVersion(1, 2, 3)
        self.assertEqual(version.major, 1)
        self.assertEqual(version.minor, 2)
        self.assertEqual(version.patch, 3)
        self.assertEqual(version.suffix, "")

    def test_constructor_with_major_minor_patch_suffix(self):
        # Test constructor with major, minor, patch, suffix
        version = InstallerVersion(1, 2, 3, "alpha")
        self.assertEqual(version.major, 1)
        self.assertEqual(version.minor, 2)
        self.assertEqual(version.patch, 3)
        self.assertEqual(version.suffix, "alpha")

    def test_name_property_without_suffix(self):
        # Test name property without suffix
        version = InstallerVersion(1, 2, 3)
        self.assertEqual(version.name, "v1.2.3")

    def test_name_property_with_suffix(self):
        # Test name property with suffix
        version = InstallerVersion("v1.2.3-alpha")
        self.assertEqual(version.name, "v1.2.3-alpha")

    def test_parse_with_invalid_version_name(self):
        # Test parse method with invalid version name
        with self.assertRaises(ValueError):
            InstallerVersion("invalid-version-name")
    
    def test_name_property_with_all_fields(self):
        # Test name property with all fields
        version = InstallerVersion(1, 2, 3, "alpha")
        self.assertEqual(version.name, "v1.2.3-alpha")

        version.major = 2
        self.assertEqual(version.name, "v2.2.3-alpha")

        version.minor = 3
        self.assertEqual(version.name, "v2.3.3-alpha")

        version.patch = 4
        self.assertEqual(version.name, "v2.3.4-alpha")

        version.suffix = "beta"
        self.assertEqual(version.name, "v2.3.4-beta")

