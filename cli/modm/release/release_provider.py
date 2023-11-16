import hashlib
import os
from pathlib import Path
import requests

from modm.release.resources_archive import ResourcesArchive
from ._config import Config
from .release_info import ReleaseInfo
from .version import Version
from .version_provider import VersionProvider
from .resources_archive import ResourcesArchiveDownloadOptions, download


class ReleaseProvider:
    """
    Enables the ability to download and cache the resources that have been added to release (that are stored in a GitHub release).
    as well as get the release info for a given version.
    """

    __instance = None
    home_dir_name = ".modm"

    def __new__(cls):
        """Singleton implementation"""
        if cls.__instance is None:
            instance = super(ReleaseProvider, cls).__new__(cls)
            instance._home_dir = instance._get_home_dir()
            instance._resources_dir = instance._get_resources_dir()
            instance._index = None

            instance._load_from_disk()
            cls.__instance = instance

        return cls.__instance

    def get(self, version: Version | str) -> ReleaseInfo:
        """
        Gets the release info for the given version.
        """
        self._load_releases_index()
        if isinstance(version, str):
            version_name = version
        else:
            version_name = version.name
        return self._index[version_name]

    def get_resources(self, version: Version | str) -> ResourcesArchive:
        """
        Gets the resources archive for the given version.
        """
        self._load_releases_index()

        if isinstance(version, str):
            version_name = version
        else:
            version_name = version.name
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

    def _load_releases_index(self):
        """
        Loads the index of releases from the index_url and caches it in memory.

        Remarks:
        The index is a json file that contains the list of all available installer versions that have been released.
        """
        # TODO: we need to implement a mechanism for saving the file to disk to reduce the amount of http requests

        if self._index is None:
            response = requests.get(Config().releases_index_url(), headers={"Accept": "application/json"})
            response.raise_for_status()
            document = response.json()

            self._index: dict[str, ReleaseInfo] = {}
            for release in document["releases"]:
                self._index[release["version"]] = ReleaseInfo.from_dict(release)

        return self._index

    def _load_from_disk(self):
        entries = os.listdir(self._resources_dir)

        # the set of resources from the resources file is a directory, named the version
        version_dirs = [Path(entry) for entry in entries if os.path.isdir(entry)]

        resources: dict[str, ResourcesArchive] = {}
        for version_dir in version_dirs:
            version = Version(version_dir.name)
            resources[version] = ResourcesArchive(version_dir, version=version)
        self._entries = resources

    def _get_resources(self, version_name) -> ResourcesArchive:
        if version_name == "latest":
            version = VersionProvider().get_latest()
        else:
            version = Version(version_name)

        if version.name in self._entries:
            return self._entries[version]
        else:
            return self._download(version)

    def _download(self, version: Version):
        """
        Downloads the resources archive tarball for the given version and extracts it to the resources directory
        into a directory named after the version.
        """
        if version.name in self._entries:
            return self._entries[version]

        if version.name not in self._index:
            raise Exception(f"Invalid version '{version.name}'. Not found in the installer release index.")

        resources_info = self._index[version.name].resources

        options = ResourcesArchiveDownloadOptions(self._resources_dir, resources_info.sha256_digest, version.name)
        resources = download(resources_info.download_url, options)
        resources.directory = self._resources_dir
        resources.extract()

        self._entries[version.name] = resources

        return resources
