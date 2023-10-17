import os
from pathlib import Path
import packaging.zip_utils as ziputils
import shutil
import tempfile
import packaging.manifest as manifest
from .manifest import ManifestInfo


class InstallerPackage:
    """
    The installer package, e.g. the installer.zip, which is a zip archive
    containing the installer's main template (and all dependencies) and the manifest file
    """

    file_name = "installer.zip"

    def __init__(self, manifest: ManifestInfo):
        self.manifest = manifest

    def create(self):
        validation_results = self.manifest.validate()
        if len(validation_results) > 0:
            raise ValueError(validation_results)

        temp_dir, templates_dir = self._get_copy_of_templates_dir()

        manifest.write(templates_dir, self.manifest)
        dest_file_path = Path(os.path.join(temp_dir, InstallerPackage.file_name))

        file = ziputils.zip_dir(templates_dir, dest_file_path)

        return file

    def unpack(self, file_path, extract_dir):
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Destination path {file_path} does not exist")

        file = Path(file_path).resolve()

        if not file.is_file():
            raise ValueError(f"Destination path {file_path} is not a file")

        shutil.unpack_archive(file, extract_dir)

    def _get_copy_of_templates_dir(self):
        source_templates_dir = Path(self.manifest.main_template).parent
        temp_dir = tempfile.mkdtemp()
        templates_dir = Path(os.path.join(temp_dir, source_templates_dir.name))

        shutil.copytree(str(source_templates_dir), templates_dir, dirs_exist_ok=True)

        return (Path(temp_dir), templates_dir)


def create_installer_package(manifest):
    """
    Creates an installer package for the given manifest.

    Args:
      manifest (ManifestInfo): instance of ManifestInfo

    Returns:
      pathlib.Path: The the installer package file as Path object.
    """
    installer_package = InstallerPackage(manifest)
    return installer_package.create()
