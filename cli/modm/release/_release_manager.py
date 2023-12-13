from pathlib import Path
import tarfile
import tempfile
from modm.installer.client_app_package import ClientAppPackage
from modm.release.version import Version
from .release_info import ReleaseInfo
from ghapi.all import GhApi
import os
import shutil
import warnings
import subprocess
from typing import Callable
from modm import _zip_utils
from ._release_repository import repository, ReleaseRepository


class ReleaseManager:
    """
    Uses: https://ghapi.fast.ai/fullapi.html
    """

    owner = "microsoft"
    repo = "commercial-marketplace-offer-deploy"

    def __init__(self, github_auth_token: str = None, out: Callable[[str], None] = None):
        self.out = out if out is not None else lambda *args: None
        self._github_auth_token = github_auth_token
        self.cwd = os.getcwd()

    def create(self, release_info: ReleaseInfo):
        version = Version(release_info.version)

        with repository(version, self._github_auth_token) as repo:
            result = self._create_resources_archive(repo.path, version)

    def _build_cli_lib(self, repo: ReleaseRepository):
        process = subprocess.Popen(["python", "-m pip install", "--upgrade build"], stdout=subprocess.PIPE, universal_newlines=True)

        self.out("Building client app.")

        while True:
            output = process.stdout.readline()
            if len(output) > 0:
                self.out("  " + output.strip())
            # Do something else
            return_code = process.poll()
            if return_code is not None:
                # Process has finished, read rest of the output
                for output in process.stdout.readlines():
                    if len(output) > 0:
                        self.out("  " + output.strip())
                break

        self.out("Creating client app package.")
    def _create_resources_archive(self, repo: ReleaseRepository) -> Path:
        out_dir = repo.path / "dist"

        client_app_package_file = self._create_client_app_package(repo.path)
        template_files = repo.template_files()

        out_file = out_dir / f"resources-{repo.version.name}.tar.gz"

        with tarfile.open(out_file, "w:gz") as tar:
            for template_file in template_files:
                tar.add(template_file, arcname=template_file.name)
            tar.add(client_app_package_file, arcname=client_app_package_file.name)

        result = {
            "file": out_file,
            "release": {
                "downloadUrl": f"{repo.release_download_url}/{out_file.name}",
                "filename": out_file.name,
                "sha256Digest": self._get_sha256_digest(out_file)
            }
        }
        return result

    def _create_client_app_package(self, clone_dir: Path) -> Path:
        """creates a clientapp.zip"""
        out_dir = clone_dir / "dist"

        if not out_dir.exists():
            out_dir.mkdir()

        csproj_file = clone_dir / "src" / "ClientApp" / "ClientApp.csproj"

        process = subprocess.Popen(["dotnet", "publish", csproj_file, "-c", "Release", "-o", out_dir], stdout=subprocess.PIPE, universal_newlines=True)

        self.out("Building client app.")

        while True:
            output = process.stdout.readline()
            if len(output) > 0:
                self.out("  " + output.strip())
            # Do something else
            return_code = process.poll()
            if return_code is not None:
                # Process has finished, read rest of the output
                for output in process.stdout.readlines():
                    if len(output) > 0:
                        self.out("  " + output.strip())
                break

        self.out("Creating client app package.")

        out_file = out_dir / ClientAppPackage.file_name
        _zip_utils.zip_dir(out_dir, out_file)

        return out_file