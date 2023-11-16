
import hashlib
from pathlib import Path
from . import _httputil as httputil
from .version import Version
import tarfile


class ResourcesArchive:
    default_file_name = "resources.tar.gz"
    file_name_format = "resources-{version_name}.tar.gz"

    """
    The archive file contained in the resources tarball that gets uploaded to the release and referenced in the index.json releases index
    """
    def __init__(self, resources_dir: Path = None, resources_file: Path = None, version: str = None):
        self._version = version

        if resources_file is not None:
            self.directory: Path = resources_file.parent
            self.file: Path = resources_file
            self.extract()
        else:
            self.directory: Path = resources_dir
            self.file: Path = None

        self.main_template = self.directory.joinpath("mainTemplate.json")
        self.view_definition = self.directory.joinpath("viewDefinition.json")
        self.create_ui_definition_step = self.directory.joinpath("createUiDefinition.json")
        self.client_app_package = self.directory.joinpath("clientapp.zip")

    @property
    def version(self) -> Version:
        return self._version

    @version.setter
    def version(self, value: str):
        self._version = value
    
    @property 
    def file_name(self):
        if self._version is not None:
            return self.file_name_format.format(version_name=self._version)
        else:
            return self.default_file_name

    def archive(self, out_dir: Path):
        """
        Tarball's the resources file and returns the tarball file path
        """
        out_file = out_dir / self.file_name

        with tarfile.open(out_file, "w:gz") as tar:
            tar.add(self.main_template, arcname="mainTemplate.json")
            tar.add(self.view_definition, arcname="viewDefinition.json")
            tar.add(self.create_ui_definition_step, arcname="createUiDefinition.json")
            tar.add(self.client_app_package, arcname="clientapp.zip")

        return out_file
    
    def extract(self):
        with tarfile.open(self.file, "r:gz") as tar:
            tar.extractall(path=self.directory)


class ResourcesArchiveDownloadOptions:
    def __init__(self, dest_dir: Path, sha256_digest: str, version: str):
        self.version = version
        self.dest_dir = dest_dir
        self.sha256_digest = sha256_digest

    def validate_sha256_digest(self, file: Path):
        actual_digest = hashlib.sha256(file.read_bytes()).hexdigest()
        if actual_digest != self.sha256_digest:
            raise Exception(
                f"Invalid resources file '{file.name}'. The expected SHA256 digest is '{self.sha256_digest}' but the actual digest is '{actual_digest}'."
            )


def download(resources_file_url: str, options: ResourcesArchiveDownloadOptions) -> ResourcesArchive:
    file_path = options.dest_dir / resources_file_url.rsplit('/', 1)[-1]

    httputil.download_file(resources_file_url, file_path)
    options.validate_sha256_digest(file_path)

    return ResourcesArchive(options.dest_dir, file_path, options.version)