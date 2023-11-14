from importlib.resources import as_file, files
import json


class Config:
    __instance = None

    def __new__(cls):
        if cls.__instance is None:
            cls.__instance = super().__new__(cls)
            cls.__instance._load_from_package_resources()
        return cls.__instance

    @classmethod
    def index_url(cls):
        return cls.get_value("index_url")

    @classmethod
    def get(cls):
        return cls.__instance._config

    @classmethod
    def get_value(cls, key):
        return cls.__instance._config.get(key)

    def _load_from_package_resources(self):
        config = files("modm.config")
        with as_file(config.joinpath("config.json")) as config_file:
            with open(config_file, 'r') as config_file:
                self._config = json.load(config_file)
