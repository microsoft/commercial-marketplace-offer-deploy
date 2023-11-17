import tempfile
from modm.marketplace.application_package_resources import ApplicationPackageResources
from modm.release.release_provider import ReleaseProvider
from modm.release.resources_archive import ResourcesArchive
from modm.release.version import Version
from pathlib import Path

from modm.release.version_provider import VersionProvider


class ApplicationPackageOptions:
    """
    Options for creating an application package.

    Args:
        installer_version (InstallerVersion | str): The version of the installer to use.
        use_vmi_reference (bool, optional): Whether to use a VMI reference of the published/released reference. Defaults to False.
        vmi_reference_id (str, optional): The ID of the VMI reference to use to override the published reference.
        out_dir (Optional[str]): The output directory for the application package.
    """

    def __init__(
        self,
        version: Version | str,
        use_vmi_reference: bool = False,
        vmi_reference_id: str = None,
        resources_archive_file: str | Path = None,
        out_dir=None,
    ) -> None:
        self._out_dir = out_dir
        self._use_vmi_reference = use_vmi_reference
        self._vmi_reference_id = vmi_reference_id

        self._set_version(version)
        self._is_version_set = version is not None
              
        self._set_resources(resources_archive_file)

        if vmi_reference_id is not None:
            self._use_vmi_reference = True

    def _set_resources(self, resources_archive_file: str):
        self._resources: ApplicationPackageResources = None
        self._resources_archive: ResourcesArchive = None

        is_file_directly_specified = resources_archive_file is not None and (isinstance(resources_archive_file, str) or isinstance(resources_archive_file, Path))

        if is_file_directly_specified:
            self._resources_archive = ResourcesArchive(resources_file=resources_archive_file, version=self.version)
            self._resources = ApplicationPackageResources(self._resources_archive)
        elif self._is_version_set:
            # If the resources file is not directly specified, and the version is set, use the released version
            provider = ReleaseProvider()
            self._resources_archive = provider.get_resources(self.version)
            self._release_reference = provider.get(self.version)

            self._resources = ApplicationPackageResources(self._resources_archive, self._release_reference)

    @property
    def resources(self) -> ApplicationPackageResources:
        return self._resources

    @property
    def vmi_reference_id(self):
        """This is the ID of the VMI reference to use to override the published reference."""
        return self._vmi_reference_id

    @property
    def use_vmi_reference(self):
        return self._use_vmi_reference

    @property
    def out_dir(self):
        """
        Returns the output directory for the packaged application.
        If the output directory is not set, a temporary directory is created and returned.
        """
        out_dir = self._out_dir if self._out_dir is not None else tempfile.mkdtemp()
        return out_dir

    def _set_version(self, version):
        if version is not None:
            if isinstance(version, str):
                if version == "latest":
                    self.version = VersionProvider().get_latest()
                else:
                    self.version = Version(version)
            else:
                self.version = version
