
import requests

from .version import Version
from ._config import Config
from . import _github_release as github


class VersionProvider:
    "" "Singleton implementation" ""
    _instance = None
    _versions = {}

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance

    def get_latest(self):
        self._load()
        latest_release = github.get_latest_release()
        return Version(latest_release["tag_name"])

    def get(self, version: str | Version):
        self._load()
        value = version
        if isinstance(value, str):
            value = Version(version)
        return self._versions.get(value.name)

    def list(self):
        self._load()
        return list(self._versions.values())

    def _load(self):
        if self._versions is None:
            response = requests.get(Config().releases_index_url(), headers={"Accept": "application/json"})
            response.raise_for_status()
            document = response.json()
            self._versions: dict[str, Version] = {}
            for release in document["releases"]:
                name = release["version"]
                self._versions[name] = Version(name)
        return self._versions
