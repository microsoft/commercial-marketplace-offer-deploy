from pathlib import Path
from . import main_template, view_definition
from modm.release.version import Version
import tarfile

from . import create_ui_definition_step


class ApplicationPackageResources:
    def __init__(self, location: Path, installer_version = None, release_reference: dict = None):
        self.installer_version = installer_version
        if release_reference is not None:
            self._release = release_reference
            self.vmi_reference_id = release_reference["vmi"]
            self.vm_offer = release_reference["offer"]

        self.location = location
        self.main_template = main_template.from_file(location.joinpath("mainTemplate.json"))
        self.view_definition = view_definition.from_file(location.joinpath("viewDefinition.json"))
        self.create_ui_definition_step = create_ui_definition_step.from_file(location.joinpath("createUiDefinition.json"))
        self.client_app_package = location.joinpath("clientapp.zip")

    @staticmethod
    def from_file(resources_file: Path):
        """this will extract the installer resources from a resources tarball file"""
        if not resources_file.exists():
            raise Exception(f"File '{resources_file}' does not exist.")
        
        with tarfile.open(resources_file, "r:gz") as tar:
            tar.extractall(path=resources_file.parent)

        return ApplicationPackageResources(resources_file.parent)
    
    def to_file(self, resources_dir: Path, version: str = None, out_dir: Path = None) -> Path:
        main_template_file = resources_dir / 'mainTemplate.json'
        view_definition_file = resources_dir / 'viewDefinition.json'
        create_ui_definition_file = resources_dir / 'createUiDefinition.json'
        client_app_file = resources_dir / 'clientapp.zip'
        
        if version is not None:
            installer_version = Version(version)
            out_file = out_dir / f"resources-{installer_version.name}.tar.gz"
        else:
            out_file = out_dir / f"resources.tar.gz"

        with tarfile.open(out_file, "w:gz") as tar:
            tar.add(main_template_file, arcname=main_template_file.name)
            tar.add(view_definition_file, arcname=view_definition_file.name)
            tar.add(create_ui_definition_file, arcname=create_ui_definition_file.name)
            tar.add(client_app_file, arcname=client_app_file.name)
        
        return out_file
