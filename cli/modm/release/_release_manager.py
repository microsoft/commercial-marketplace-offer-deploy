from pathlib import Path
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

class ReleaseManager:
    """
    Uses: https://ghapi.fast.ai/fullapi.html
    """

    owner = "microsoft"
    repo = "commercial-marketplace-offer-deploy"

    def __init__(self, github_auth_token: str = None, out: Callable[[str], None] = None):
        self.out = out if out is not None else lambda *args: None
        self.api: GhApi = self._get_api_client(github_auth_token)
        self.cwd = os.getcwd()

    def create(self, release_info: ReleaseInfo):
        version = Version(release_info.version)

        if self._release_exists(version):
            raise ValueError(f"Release {version.name} already exists")

        self._create_tag(version)
        clone_dir = self._create_shallow_clone(version)

        resources_archive = self._create_resources_archive(clone_dir)

        # 5. build the wheel of the python library for the version
        wheel_file = "/path/to/wheel/file"
        os.system("python setup.py bdist_wheel")

        # 6. create a draft GH release for the tagged version
        release = self.github.repos.create_release(tag_name=version.name, draft=True)

        # 7. upload the files
        #  - the wheel
        self.github.repos.upload_release_asset(release_id=release.id, name="wheel", asset=wheel_file)
        
        #  - the resources archive
        self.github.repos.upload_release_asset(release_id=release.id, name="resources", asset=resources_archive)

    def _create_resources_archive(self, clone_dir: Path) -> Path:
        cwd = Path(current_working_dir)
        out_dir = _resolve_path(cwd, out_dir)

        templates_dir = clone_dir / "templates"
        main_template_file = templates_dir / "mainTemplate.json"
        view_definition_file = templates_dir/ "viewDefinition.json"
        create_ui_definition_file = templates_dir / "createUiDefinition.json"
        client_app_package_file = self._create_client_app_package(clone_dir)

        if version is not None:
            installer_version = Version(version)
            out_file = out_dir / f"resources-{installer_version.name}.tar.gz"
        else:
            out_file = out_dir / f"resources.tar.gz"

        click.echo(f"Creating tarball {out_file}...")

        with tarfile.open(out_file, "w:gz") as tar:
            tar.add(main_template_file, arcname=main_template_file.name)
            tar.add(view_definition_file, arcname=view_definition_file.name)
            tar.add(create_ui_definition_file, arcname=create_ui_definition_file.name)
            tar.add(client_app_file, arcname=client_app_file.name)

        click.echo(f"resources '{out_file.name}' created.")

        if version is not None:
            result = {}
            result[
                "downloadUrl"
            ] = f"https://github.com/microsoft/commercial-marketplace-offer-deploy/releases/download/{installer_version.name}/{out_file.name}"
            result["filename"] = out_file.name
            result["sha256Digest"] = _get_sha256_digest(out_file)
        else:
            result = {}
            result["filename"] = str(out_file)
            result["sha256Digest"] = _get_sha256_digest(out_file)

        click.echo(json.dumps(result, indent=2))

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

    def _get_api_client(self, github_auth_token: str = None):
        """
        Returns a GitHub API client
        If github_auth_token is None, the client will be attempt to authenticate using the standard GitHub environment
        variables (which is the default)
        """

        # the GitHub API client only warns if authentication fails. We need to catch this and exit if it happens
        warnings.filterwarnings("error")
        api = None

        try:
            api = GhApi(owner="microsoft", repo="commercial-marketplace-offer-deploy", token=github_auth_token, authenticate=True)
        except UserWarning as e:
            if "unauthenticated" in str(e):
                err = RuntimeError("Authentication failed")

        warnings.resetwarnings()

        # try to get something with the authenticated user. if it fails, raise an error
        try:
            api.repos.list_for_authenticated_user()
        except Exception as e:
            err = e

        if err is not None:
            raise err
        return api

    def _release_exists(self, version: Version):
        """Returns True if a release with the given version already exists"""
        releases = self.api.repos.list_releases()
        for release in releases:
            if release.tag_name == version:
                return True
        return False

    def _create_tag(self, version: Version):
        """Creates a tag for the given version"""
        main_ref = self.api.git.get_ref("heads/main")
        tag = self.api.git.create_tag(tag=version.name, message=version.name, object=main_ref.object.sha, type="commit")
        self.api.git.create_ref(ref=f"refs/tags/{version.name}", sha=tag.sha)

    def _create_shallow_clone(self, version: Version):
        """Creates a shallow clone of the given version"""
        temp_dir = tempfile.mkdtemp()
        os.chdir(temp_dir)
        repo_url = f"https://github.com/{ReleaseManager.owner}/commercial-marketplace-offer-deploy.git"
        os.system(f"git clone --depth 1 --branch {version.name} {repo_url} .")

        return Path(temp_dir)