
from importlib.resources import as_file, files
from pathlib import Path
from urllib.parse import urlparse
import urllib.request


class ClientAppPackage:
    """
    A class that represents a client app package which is the frontend ClientApp
    that will be deployed to app services.
    """
    file_name = "clientapp.zip"

    def __init__(self, file):
        self.file = file

    @property
    def path(self):
        return self.file

    @property
    def name(self):
        return self.file.name
    
    @staticmethod
    def from_uri(url: str, dest_dir: str):
        parsed_url = urlparse(url)
        file = Path(dest_dir, parsed_url.basename).resolve()
        _, out_path = urllib.request.urlretrieve(url, file)

        return ClientAppPackage(Path(out_path))

    @staticmethod
    def from_resource():
        resource_files = files("resources")
        with as_file(resource_files.joinpath(ClientAppPackage.file_name)) as resource_file:
            package = ClientAppPackage(resource_file)
            return package