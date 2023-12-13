import os
from pathlib import Path
import tempfile
import warnings
from ghapi.all import GhApi
from modm.release.version import Version


class ReleaseRepository:
    host = "github.com"
    owner = "microsoft"
    repo = "commercial-marketplace-offer-deploy"

    def __init__(self, version: Version, clone_path: Path = tempfile.mkdtemp(), github_auth_token: str = None):
        self._path = clone_path
        self._version = version
        self._client: GhApi = self._get_api_client(github_auth_token)
        
    def __enter__(self):
        self.cwd = os.getcwd()
        os.chdir(self.path)

    def __exit__(self, exc_type, exc_val, exc_tb):
        os.chdir(self.cwd)

    @property
    def version(self) -> Version:
        return self._version

    @property
    def path(self):
        return self._path

    @property
    def release_download_url(self):
        return f"https://{self.host}/{self.owner}/{self.repo}/releases/download/{self._version.name}"

    def release_exists(self):
        """Returns True if a release with the given version already exists"""
        releases = self._client.repos.list_releases()
        for release in releases:
            if release.tag_name == self._version:
                return True
        return False

    def create_release(self, assets: list[Path]):
        if self.release_exists(self._version):
            raise ValueError(f"Release {self._version.name} already exists")
        release = self._client.repos.create_release(tag_name=self._version.name, draft=True)

        for asset in assets:
            self._client.repos.upload_release_asset(release_id=release.id, name=asset.name, asset=asset)

        return release

    def template_files(self):
        """Returns list of the template files"""
        return list(self.path.glob("templates/*.json"))

    def clone(self):
        """Creates a shallow clone of the given version"""
        repo_url = f"https://{self.host}/{self.owner}/{self.repo}.git"
        os.system(f"git clone --depth 1 --branch {self.version.name} {repo_url} .")

    def tag(self):
        """Creates a tag for the given version"""
        main_ref = self._client.git.get_ref("heads/main")
        tag = self._client.git.create_tag(tag=self._version.name, message=self._version.name, object=main_ref.object.sha, type="commit")
        self._client.git.create_ref(ref=f"refs/tags/{self._version.name}", sha=tag.sha)

    def _get_api_client(self, github_auth_token: str = None):
        """
        Returns a GitHub API client
        If github_auth_token is None, the client will be attempt to authenticate using the standard GitHub environment
        variables (which is the default)
        """

        # the GitHub API client only warns if authentication fails. We need to catch this and exit if it happens
        warnings.filterwarnings("error")
        client = None

        try:
            client = GhApi(owner=self.owner, repo=self.repo, token=github_auth_token, authenticate=True)
        except UserWarning as e:
            if "unauthenticated" in str(e):
                err = RuntimeError("Authentication failed")

        warnings.resetwarnings()

        # try to get something with the authenticated user. if it fails, raise an error
        try:
            client.repos.list_for_authenticated_user()
        except Exception as e:
            err = e

        if err is not None:
            raise err
        return client

def repository(version: Version, github_auth_token: str = None) -> ReleaseRepository:
    repo =  ReleaseRepository(version, github_auth_token=github_auth_token)
    if repo.release_exists():
        raise ValueError(f"Release {version.name} already exists")

    repo.tag()
    repo.clone()
    return repo