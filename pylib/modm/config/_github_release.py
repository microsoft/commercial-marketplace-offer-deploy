import requests

_latest_release_url = "https://api.github.com/repos/microsoft/commercial-marketplace-offer-deploy/releases/latest"
_release_by_tag_url_format = "https://api.github.com/repos/microsoft/commercial-marketplace-offer-deploy/releases/tags/{}"

# Further information can be found at https://docs.github.com/en/rest/reference/repos#releases


def get_latest_release():
    """Get the latest version of the installer from GitHub and return the InstallerVersion instance"""
    return _make_request(_latest_release_url)


def get_release_by_tag_name(name):
    """
    The tag name is the equivalent of the full version name including the "v" prefix; example: v1.0.0
    """
    return _make_request(_release_by_tag_url_format.format(name))


def _make_request(url):
    headers = {"Accept": "application/vnd.github.v3+json", "X-GitHub-Api-Version": "2022-11-28"}
    response = requests.get(url, headers=headers)
    response.raise_for_status()
    return response.json()
