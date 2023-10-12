
# create a class called InstallerPackage
# it has a static field on it called file_name
import os
from pathlib import Path
from packaging import manifest_file
from packaging.deployment_type import DeploymentType
from packaging.manifest import ManifestInfo
import packaging.zip_utils as ziputils
import shutil
import tempfile

class InstallerPackage:
  file_name = 'installer.pkg'

  def __init__(self, manifest: ManifestInfo):
    self.manifest = manifest

  def create(self):
    validation_results = self.manifest.validate()
    if len(validation_results) > 0:
      raise ValueError(validation_results)

    temp_dir, templates_dir = self._get_copy_of_templates_dir()

    manifest_file.write(templates_dir, self.manifest)
    dest_file_path = Path(os.path.join(temp_dir, InstallerPackage.file_name))

    file = ziputils.zip_dir(templates_dir, dest_file_path)

    return file, temp_dir

  def unpack(self, file_path, extract_dir):
    if not os.path.exists(file_path):
      raise FileNotFoundError(f'Destination path {file_path} does not exist')

    file = Path(file_path).resolve()

    if not file.is_file():
      raise ValueError(f'Destination path {file_path} is not a file')
    
    if file.suffix != '.pkg':
      raise ValueError(f'Destination file {file_path} must have a .pkg extension')
    
    archive_file = file.with_suffix('.zip')
    shutil.copyfile(str(file), archive_file)

    shutil.unpack_archive(archive_file, extract_dir)
  
  def _get_copy_of_templates_dir(self):
    source_templates_dir = Path(self.manifest.main_template).parent
    print(source_templates_dir)
    temp_dir =  tempfile.mkdtemp()
    templates_dir = Path(os.path.join(temp_dir, source_templates_dir.name))

    shutil.copytree(str(source_templates_dir), templates_dir, dirs_exist_ok=True)

    return (Path(temp_dir), templates_dir)