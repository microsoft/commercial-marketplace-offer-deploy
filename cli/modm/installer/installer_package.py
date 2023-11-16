import os
from pathlib import Path
import modm._zip_utils as ziputils
import shutil
import tempfile
from modm.arm.bicep_template_compiler import BicepTemplateCompiler

from modm.installer.solution_template_type import SolutionTemplateType

from .installer_package_result import InstallerPackageResult
from .manifest import ManifestInfo, write_manifest

class InstallerPackage:
    """
    The installer package, e.g. the installer.zip, which is a zip archive
    containing the installer's main template (and all dependencies) and the manifest file
    """

    file_name = "installer.zip"

    def __init__(self, manifest: ManifestInfo):
        self.manifest = manifest

    def create(self) -> InstallerPackageResult:
        validation_results = self.manifest.validate()
        if len(validation_results) > 0:
            raise ValueError(validation_results)

        parent_working_dir, templates_dir = self.get_solution_template_dir()

        self._write_manifest(templates_dir)

        installer_package_file_path = Path(os.path.join(parent_working_dir, InstallerPackage.file_name))
        installer_package_file = ziputils.zip_dir(templates_dir, installer_package_file_path)

        self.manifest.dispose()
        
        return InstallerPackageResult(installer_package_file)

    def unpack(self, file_path, extract_dir):
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Destination path {file_path} does not exist")

        file = Path(file_path).resolve()

        if not file.is_file():
            raise ValueError(f"Destination path {file_path} is not a file")

        shutil.unpack_archive(file, extract_dir)

    def get_solution_template_dir(self):
        """
        Gets the solution template from the solution template directory to a temporary directory
        and returns the temp parent directory and the templates directory.
        """
        src_templates_dir = Path(self.manifest.solution_template).parent
        dest_dir = Path(tempfile.mkdtemp())

        new_templates_dir = Path(os.path.join(dest_dir, src_templates_dir.name))
        new_templates_dir.mkdir()

        self._copy_dir(src_templates_dir, new_templates_dir)

        if self.manifest.has_bicep_source:
            bicep_dir = new_templates_dir / ".bicep"
            self._copy_dir(self.manifest.bicep_templates_dir, bicep_dir)

        return (dest_dir, new_templates_dir)

    def _write_manifest(self, templates_dir):
        write_manifest(templates_dir, self.manifest)
    
    def _copy_dir(self, src_dir: Path, dest_dir):
        shutil.copytree(str(src_dir), str(dest_dir), dirs_exist_ok=True)

def create_installer_package(manifest) -> InstallerPackageResult:
    """
    Creates an installer package for the given manifest.

    Args:
      manifest (ManifestInfo): instance of ManifestInfo

    Returns:
      pathlib.Path: The the installer package file as Path object.
    """
    installer_package = InstallerPackage(manifest)
    return installer_package.create()
