
from importlib.resources import as_file, files
from pathlib import Path
from urllib.parse import urlparse
import urllib.request


class FunctionAppPackage:
    file_name = "function.zip"

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

        return FunctionAppPackage(Path(out_path))

    @staticmethod
    def from_resource():
        resource_files = files("resources")
        with as_file(resource_files.joinpath(FunctionAppPackage.file_name)) as resource_file:
            package = FunctionAppPackage(resource_file)
            return package