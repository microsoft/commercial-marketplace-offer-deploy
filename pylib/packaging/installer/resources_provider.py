import hashlib
import os
from pathlib import Path

import requests
from packaging.config import Config
from packaging.installer.resources import InstallerResources, ResourcesInfo
from packaging.installer.version import InstallerVersion, InstallerVersionProvider
from . import _httputil as httputil
import tarfile


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
            document = response.json()
            self._index = {}
            for release in document["releases"]:
                self._index[release["version"]] = release
        return self._index

    def _load_from_disk(self):
        entries = os.listdir(self._resources_dir)
        version_dirs = [Path(entry) for entry in entries if os.path.isdir(entry)]

        resources = {}
        for version_dir in version_dirs:
            version = InstallerVersion(version_dir.name)
            resources[version] = InstallerResources(version_dir, version, self._index[version.name]["reference"])
        self._entries = resources

    def _get_resources(self, version_name):
        if version_name == "latest":
            version = InstallerVersionProvider().get_latest()
        else:
            version = InstallerVersion(version_name)

        if version in self._entries:
            return self._entries[version]
        else:
            return self._download(version)

    def _download(self, version: InstallerVersion):
        if version in self._entries:
            return self._entries[version]

        if version.name not in self._index:
            raise Exception(f"Invalid version '{version.name}'. Not found in the installer release index.")

        resources = self._fetch_resources(self._index[version.name], version)
        self._entries[version] = resources

        return resources

    def _fetch_resources(self, release: dict, version: InstallerVersion) -> InstallerResources:
        resources: ResourcesInfo = ResourcesInfo.from_dict(release["resources"])
        file_path = self._resources_dir / resources.filename

        httputil.download_file(resources.download_url, file_path)
        self._validate_sha256_digest(file_path, resources.sha256_digest)

        extract_path = self._get_resources_dir() / version.name
        extract_path.mkdir(exist_ok=True)

        with tarfile.open(file_path, "r:gz") as tar:
            tar.extractall(path=extract_path)

        return InstallerResources(extract_path, version, release["reference"])

    def _validate_sha256_digest(self, file: Path, expected_digest):
        actual_digest = hashlib.sha256(file.read_bytes()).hexdigest()
        if actual_digest != expected_digest:
            raise Exception(
                f"Invalid resources file '{file.name}'. The expected SHA256 digest is '{expected_digest}' but the actual digest is '{actual_digest}'."
            )
