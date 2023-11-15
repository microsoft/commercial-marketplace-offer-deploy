import tempfile
from modm.release.version import Version
from pathlib import Path

from modm.release.version_provider import VersionProvider


class ApplicationPackageOptions:
    """
    Options for creating an application package.

    Args:
        installer_version (InstallerVersion | str): The version of the installer to use.
        vmi_reference (bool, optional): Whether to use a VMI reference of the published/released reference. Defaults to False.
        vmi_reference_id (str, optional): The ID of the VMI reference to use to override the published reference.
        out_dir (Optional[str]): The output directory for the application package.
    """

    def __init__(
        self,
        version: Version | str,
        vmi_reference: bool = False,
        vmi_reference_id: str = None,
        resources_file: str | Path = None,
        out_dir=None,
    ) -> None:
        self._out_dir = out_dir
        self._use_vmi_reference = vmi_reference
        self._vmi_reference_id = vmi_reference_id

        if isinstance(resources_file, str):
            self._resources_file = Path(resources_file)
        else:
            self._resources_file = resources_file

        if version is not None:
            if isinstance(version, str):
                if version == "latest":
                    self.version = VersionProvider().get_latest()
                else:
                    self.version = Version(version)
            else:
                self.version = version

        if vmi_reference_id is not None:
            self._use_vmi_reference = True

    @property
    def resources_file(self):
        return self._resources_file

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
