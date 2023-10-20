import hashlib
import os
from pathlib import Path

import requests
from packaging.config import Config
from packaging.installer import main_template, view_definition
from packaging.installer.version import InstallerVersion
from . import _github_release as github
from . import _httputil as httputil
from msrest.serialization import Model
import tarfile


class ResourcesInfo(Model):
    _attribute_map = {
        "download_url": {"key": "downloadUrl", "type": "str"},
        "filename": {"key": "filename", "type": "str"},
        "sha256_digest": {"key": "sha256Digest", "type": "str"},
    }

    def __init__(self, download_url=None, filename=None, sha256_digest=None):
        super(ResourcesInfo, self).__init__()
        self.download_url = download_url
        self.filename = filename
        self.sha256_digest = sha256_digest


class InstallerResources:
    def __init__(self, location: Path, installer_version):
        self.installer_version = installer_version
        self.location = location
        self.main_template = main_template.from_file(location.joinpath("templates/mainTemplate.json"))
        self.view_definition = view_definition.from_file(location.joinpath("templates/viewDefinition.json"))
        self.function_app_package = location.joinpath("function.zip")
        self.vmi_reference_id = None
        self.vm_offer = None


class InstallerResourcesProvider:
    __instance = None
    home_dir_name = ".modm"

    def __new__(cls):
        """Singleton implementation"""
        if cls.__instance is None:
            instance = super(InstallerResourcesProvider, cls).__new__(cls)
            instance._home_dir = instance._get_home_dir()
            instance._resources_dir = instance._get_resources_dir()
            instance._index = None

            instance._load_from_disk()
            cls.__instance = instance

        return cls.__instance

    def get(self, version_name: str):
        self._load_index()
        return self._get_resources(version_name)

    def _get_home_dir(self):
        home_dir = Path.home().joinpath(self.home_dir_name)
        if not home_dir.exists():
            home_dir.mkdir()
        return home_dir

    def _get_resources_dir(self) -> Path:
        resources_dir = self._home_dir.joinpath("resources")
        if not resources_dir.exists():
            resources_dir.mkdir()
        return resources_dir

    def _load_index(self):
        """
        Loads the index from the index_url and caches it in memory.

        Remarks:
        The index is a json file that contains the list of all available installer versions that have been released.
        """
        # TODO: we need to implement a mechanism for saving the file to disk to reduce the amount of http requests

        if self._index is None:
            response = requests.get(Config().index_url(), headers={"Accept": "application/json"})
            response.raise_for_status()
            self._index = response.json()
        return self._index

    def _load_from_disk(self):
        entries = os.listdir(self._resources_dir)
        version_dirs = [Path(entry) for entry in entries if os.path.isdir(entry)]

        resources = {}
        for version_dir in version_dirs:
            version = InstallerVersion(version_dir.name)
            resources[version] = InstallerResources(version_dir, version)
        self._entries = resources

    def _get_resources(self, version_name):
        if version_name == "latest":
            version = self._get_latest_version()
        else:
            version = InstallerVersion(version_name)

        if version in self._entries:
            return self._entries[version]
        else:
            return self._download(version)

    def _download(self, version: InstallerVersion):
        if version in self._entries:
            return self._entries[version]

        releases = self._index["releases"]
        release = None
        for r in releases:
            if release["version"] == version.name:
                release = r
                break
        if release is None:
            raise Exception(f"Invalid version '{version.name}'. Not found in the official index.")

        resources = self._get_resources(release["resources"], version)
        self._entries[version] = resources

        return resources

    def _get_latest_version(self):
        latest_release = github.get_latest_release()
        version_name = latest_release["tag_name"]
        return InstallerVersion(version_name)

    def _get_resources(self, resources_dict: dict, version: InstallerVersion) -> InstallerResources:
        resources: ResourcesInfo = ResourcesInfo.from_dict(resources_dict)
        file_path = self._resources_dir / resources.filename

        httputil.download_file(resources.download_url, file_path)
        self._validate_sha256_digest(file_path, resources.sha256_digest)

        extract_path = self._get_resources_dir() / version.name
        extract_path.mkdir(exist_ok=True)

        with tarfile.open(file_path, "r:gz") as tar:
            tar.extractall(path=extract_path)

        return InstallerResources(extract_path, version)

    def _validate_sha256_digest(self, file: Path, expected_digest):
        actual_digest = hashlib.sha256(file.read_bytes()).hexdigest()
        if actual_digest != expected_digest:
            raise Exception(
                f"Invalid resources file '{file.name}'. The expected SHA256 digest is '{expected_digest}' but the actual digest is '{actual_digest}'."
            )
