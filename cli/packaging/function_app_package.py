
from pathlib import Path
from urllib.parse import urlparse
import urllib.request


class FunctionAppPackage:
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