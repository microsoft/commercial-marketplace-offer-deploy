from modm.release.release_info import ReferenceInfo
from modm.release.resources_archive import ResourcesArchive
from . import main_template, view_definition
from modm.release.version import Version
from . import create_ui_definition_step


class ApplicationPackageResources:
    def __init__(self, resources_archive: ResourcesArchive, release_reference: ReferenceInfo = None):
        self.version : Version = resources_archive.version
        self.release_reference = release_reference

        self.main_template = main_template.from_file(resources_archive.main_template)
        self.view_definition = view_definition.from_file(resources_archive.view_definition)
        self.create_ui_definition_step = create_ui_definition_step.from_file(resources_archive.create_ui_definition_step)
        self.client_app_package = resources_archive.client_app_package
