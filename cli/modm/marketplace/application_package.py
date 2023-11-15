from pathlib import Path
from zipfile import ZipFile
from modm.release.release_info import ReferenceInfo
from modm.release.release_provider import ReleaseProvider

from modm.release.resources_archive import ResourcesArchive
from .application_package_result import ApplicationPackageResult
from .application_packaging_options import ApplicationPackageOptions
from .application_package_info import ApplicationPackageInfo
from ._resources import InstallerResources
from ._resources_provider import ReleaseProvider
from .main_template_finalizer import MainTemplateFinalizer
from modm.installer import InstallerPackageResult, create_installer_package


MAIN_TEMPLATE_FILE_NAME = "mainTemplate.json"
CREATE_UI_DEFINITION_FILE_NAME = "createUiDefinition.json"
VIEW_DEFINITION_FILE_NAME = "viewDefinition.json"


class ApplicationPackage:
    """
    Represents the app package, e.g. the app.zip
    The installer package (installer.zip) will reside directly in the app.zip next to
    the installer's mainTemplate.json and createUiDefinition.json, respectively

    Example Output (structure):
    - app.zip
    - mainTemplate.json
    - createUiDefinition.json
    - viewDefinition.json
    - clientapp.zip
    - installer.zip
        - manifest.json
        - main.ts (the installer's main template)
        - modules
        - <modules>
    """
    
    file_name = "app.zip"

    def __init__(self, releaseProvider: ReleaseProvider) -> None:
        self._releaseProvider: ReleaseProvider = releaseProvider

    def create(self, info: ApplicationPackageInfo, options: ApplicationPackageOptions) -> ApplicationPackageResult:
        """
        Creates an application package based on the current manifest and UI definition.

        Args:
            info (ApplicationPackageInfo): The application package info.
            If not specified, the package will be created in a randomly generated temp directory.

        Returns:
            CreateApplicationPackageResult: A result object containing the validation results and the path to the
            created application package file.
        """
        validation_results = info.validate()

        if len(validation_results) > 0:
            return ApplicationPackageResult(validation_results=validation_results)

        installer_package = create_installer_package(info.manifest)

        self._finalize_main_template(info, installer_package, options)
        self._finalize_view_definition(info, options)
        self._finalize_create_ui_definition(info, options)

        result = ApplicationPackageResult(installer_package=installer_package, client_app_name=self.main_template.client_app_name)
        result.file = self._zip(installer_package, options)

        if result.file is None or not result.file.exists():
            result.validation_results.append(Exception("Failed to create application package"))
            return result

        return result

    def _finalize_view_definition(self, info: ApplicationPackageInfo, options: ApplicationPackageOptions):
        view_definition = self._get_resources(options).view_definition
        view_definition.add_input("dashboardUrl", self.main_template.dashboard_url)
        view_definition.add_input("offerName", info.manifest.offer.name)
        view_definition.add_input("offerDescription", info.manifest.offer.description)

        self.view_definition = view_definition

    def _finalize_main_template(
        self, info: ApplicationPackageInfo, installer_package: InstallerPackageResult, options: ApplicationPackageOptions
    ):
        """the app package's main template, NOT the installer.zip container the solution template"""
        resources_archive = self._get_resources(options)

        finalizer = MainTemplateFinalizer(resources_archive.main_template)
        self.main_template = finalizer.finalize(
            template_parameters=info.template_parameters,
            reference_info=self._get_reference(options),
            installer_package=installer_package,
            use_vmi_reference=options.use_vmi_reference,
            vmi_reference_id=options.vmi_reference_id,
        )

    def _finalize_create_ui_definition(self, info: ApplicationPackageInfo, options: ApplicationPackageOptions):
        create_ui_definition = info.create_ui_definition
        create_ui_definition_step = self._get_resources(options).create_ui_definition_step
        create_ui_definition_step.append_to(create_ui_definition)

        self.create_ui_definition = create_ui_definition

    def _zip(self, installer_package, options: ApplicationPackageOptions) -> Path:
        resources_archive = self._get_resources(options)
        file = Path(options.out_dir).joinpath(self.file_name)

        with ZipFile(file, "w") as zip_file:
            zip_file.write(resources_archive.client_app_package, resources_archive.client_app_package.name)
            zip_file.write(installer_package.path, installer_package.name)
            zip_file.writestr(MAIN_TEMPLATE_FILE_NAME, self.main_template.to_json())
            zip_file.writestr(VIEW_DEFINITION_FILE_NAME, self.view_definition.to_json())
            zip_file.writestr(CREATE_UI_DEFINITION_FILE_NAME, self.create_ui_definition.to_json())
            zip_file.close()

        return file

    def _get_reference(self, options: ApplicationPackageOptions) -> ReferenceInfo:
        """Gets the reference information for the released version"""
        release = self._releaseProvider.get(options.installer_version)
        
        if release is not None:
            return release.reference
        return None
    
    def _get_reference(self, options: ApplicationPackageOptions) -> ResourcesArchive:
        return self._releaseProvider.get_resources(options.installer_version)

def new_application_package() -> ApplicationPackage:
    return ApplicationPackage(ReleaseProvider())
