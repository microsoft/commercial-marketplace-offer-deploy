import copy
import os
import jsonschema
from pathlib import Path
from packaging.manifest import ManifestInfo


class ManifestFile:
  file_name = "manifest.json"

  @staticmethod
  def write(dest_path, manifest : ManifestInfo):
    manifest_copy = copy.deepcopy(manifest)
    manifest_copy.main_template = os.path.basename(manifest.main_template)

    json = manifest_copy.to_json()
    file_path = Path(os.path.join(dest_path, ManifestFile.file_name)).resolve()

    with open(file_path, 'w') as f:
      f.write(json)

  def validate(self, schema):
    try:
      jsonschema.validate(self.read(), schema)
    except jsonschema.exceptions.ValidationError as err:
      raise err


def write(dest_path, manifest : ManifestInfo):
  ManifestFile.write(dest_path, manifest)
