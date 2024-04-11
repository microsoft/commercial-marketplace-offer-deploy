import requests
from ._config import Config

# Further information can be found at https://docs.github.com/en/rest/reference/repos#releases

def get_latest_release():
    """Get the latest version of the installer from GitHub and return the InstallerVersion instance"""
    return _make_github_api_request(Config().latest_release_url())


def get_release_by_tag_name(name):
    """
    The tag name is the equivalent of the full version name including the "v" prefix; example: v1.0.0
    """
    return _make_github_api_request(Config().release_by_tag_url_format.format(name))


def _make_github_api_request(url):
    headers = {"Accept": "application/vnd.github.v3+json", "X-GitHub-Api-Version": "2022-11-28"}
    response = requests.get(url, headers=headers)
    response.raise_for_status()
    return response.json()
