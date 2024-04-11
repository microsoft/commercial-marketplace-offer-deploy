from pathlib import Path
from zipfile import ZipFile
from modm.marketplace.application_package_resources import ApplicationPackageResources
from modm.release.release_info import ReferenceInfo
from modm.release.release_provider import ReleaseProvider
from .application_package_result import ApplicationPackageResult
from .application_packaging_options import ApplicationPackageOptions
from .application_package_info import ApplicationPackageInfo
from modm.release.release_provider import ReleaseProvider
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

    def __init__(self, info: ApplicationPackageInfo) -> None:
        """
        Args:
            info (ApplicationPackageInfo): The application package info.
        """
        self.info = info

    def create(self, options: ApplicationPackageOptions) -> ApplicationPackageResult:
        """
        Creates an application package based on the current manifest and UI definition.

        Args:
            options (ApplicationPackageOptions): The application packaging options

        Returns:
            ApplicationPackageResult: A result object containing the validation results and the path to the
            created application package file.
        """
        validation_results = self.info.validate()

        if len(validation_results) > 0:
            return ApplicationPackageResult(validation_results=validation_results)

        installer_package = create_installer_package(self.info.manifest)

        self._finalize_main_template(installer_package, options)
        self._finalize_view_definition(options)
        self._finalize_create_ui_definition(options)

        result = ApplicationPackageResult(installer_package=installer_package)
        result.file = self._zip(installer_package, options)

        if result.file is None or not result.file.exists():
            result.validation_results.append(Exception("Failed to create application package"))
            return result

        return result

    def _finalize_view_definition(self, options: ApplicationPackageOptions):
        view_definition = options.resources.view_definition
        view_definition.add_input("offerName", self.info.manifest.offer.name)
        view_definition.add_input("offerDescription", self.info.manifest.offer.description)

        self.view_definition = view_definition

    def _finalize_main_template(self, installer_package: InstallerPackageResult, options: ApplicationPackageOptions):
        """the app package's main template, NOT the installer.zip container the solution template"""
        finalizer = MainTemplateFinalizer(options.resources.main_template)
        self.main_template = finalizer.finalize(
            template_parameters=self.info.template_parameters,
            release_info=options.resources.release_reference,
            installer_package=installer_package,
            use_vmi_reference=options.use_vmi_reference,
            vmi_reference_id=options.vmi_reference_id,
        )

    def _finalize_create_ui_definition(self, options: ApplicationPackageOptions):
        create_ui_definition = self.info.create_ui_definition
        create_ui_definition_step = options.resources.create_ui_definition_step
        create_ui_definition_step.append_to(create_ui_definition)

        self.create_ui_definition = create_ui_definition

    def _zip(self, installer_package, options: ApplicationPackageOptions) -> Path:
        file = Path(options.out_dir).joinpath(self.file_name)

        with ZipFile(file, "w") as zip_file:
            zip_file.write(options.resources.client_app_package, options.resources.client_app_package.name)
            zip_file.write(installer_package.path, installer_package.name)
            zip_file.writestr(MAIN_TEMPLATE_FILE_NAME, self.main_template.to_json())
            zip_file.writestr(VIEW_DEFINITION_FILE_NAME, self.view_definition.to_json())
            zip_file.writestr(CREATE_UI_DEFINITION_FILE_NAME, self.create_ui_definition.to_json())
            zip_file.close()

        return file
