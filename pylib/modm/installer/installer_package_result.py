import hashlib
from pathlib import Path


class InstallerPackageResult:
    """
    The result of creating an installer package
    """

    def __init__(self, file: Path):
        self.file = file
        self._hash = None

    @property
    def name(self):
        return self.file.name

    @property
    def path(self) -> Path:
        return self.file

    @property
    def hash(self):
        if self._hash is None:
            self._hash = self._compute_sha256(self.file)
        return self._hash

    def _compute_sha256(self, file_name):
        hash_sha256 = hashlib.sha256()
        with open(file_name, "rb") as f:
            for chunk in iter(lambda: f.read(4096), b""):
                hash_sha256.update(chunk)
        return hash_sha256.hexdigest()

    def __str__(self):
        return self.file