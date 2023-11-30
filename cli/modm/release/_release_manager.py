
from pathlib import Path
import tempfile
from modm.release.version import Version
from .release_info import ReleaseInfo
from ghapi.all import GhApi
from fastcore.net import HTTP4xxClientError
import os
import shutil
import warnings

class ReleaseManager:
    """
    Uses: https://ghapi.fast.ai/fullapi.html
    """
    owner = "microsoft"
    repo = "commercial-marketplace-offer-deploy"

    def __init__(self):
        self.api: GhApi = self._get_api_client()
        self.cwd = os.getcwd()

    def _get_api_client(self):
        """Returns a GitHub API client"""

        # the GitHub API client only warns if authentication fails. We need to catch this and exit if it happens
        warnings.filterwarnings("error")
        api = None

        try:
            api = GhApi(owner='microsoft', repo="commercial-marketplace-offer-deploy", authenticate=True)
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
        main_ref = self.api.git.get_ref('heads/main')
        tag = self.api.git.create_tag(tag=version.name, message=version.name, object=main_ref.object.sha, type="commit")
        self.api.git.create_ref(ref=f"refs/tags/{version.name}", sha=tag.sha)

    def _create_shallow_clone(self, version: Version):
        """Creates a shallow clone of the given version"""
        temp_dir = tempfile.mkdtemp()
        os.chdir(temp_dir)
        repo_url = f"https://github.com/microsoft/commercial-marketplace-offer-deploy.git"
        os.system(f"git clone --depth 1 --branch {version.name} {repo_url} .")

        return Path(temp_dir)

    def create(self, release_info: ReleaseInfo):
        version = Version(release_info.version)
        
        if self._release_exists(version):
            raise ValueError(f"Release {version.name} already exists")

        self._create_tag(version)
        clone_dir = self._create_shallow_clone(version)
        
        # 4. build the resources archive for the version
        resources_archive = "/path/to/resources/archive"
        shutil.make_archive(resources_archive, "zip", "/path/to/resources")

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

