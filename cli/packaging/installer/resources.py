

class InstallerResources:
    latest_release_url = "https://api.github.com/repos/microsoft/commercial-marketplace-deploy-tool/releases/latest"
    release_url_format = "https://api.github.com/repos/microsoft/commercial-marketplace-deploy-tool/releases/tags/{}"

    def __init__(self, installer_version):
        self.installer_version = installer_version
