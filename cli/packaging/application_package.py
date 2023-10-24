from pathlib import Path
import tempfile
from zipfile import ZipFile
from packaging.installer import ManifestInfo, CreateInstallerPackageResult, create_installer_package
from packaging.azure import CreateUiDefinition
from importlib.resources import files, as_file
from msrest.serialization import Model
from packaging.installer import MainTemplateFinalizer
from packaging.installer.resources import InstallerResourcesProvider
from packaging.installer.version import InstallerVersion


MAIN_TEMPLATE_FILE_NAME = "mainTemplate.json"
CREATE_UI_DEFINITION_FILE_NAME = "createUiDefinition.json"
VIEW_DEFINITION_FILE_NAME = "viewDefinition.json"


class CreateApplicationPackageResult(Model):
    _attribute_map = {
        "file": {"key": "file", "type": "str"},
        "validation_results": {"key": "validationResults", "type": "[object]"},
    }

    def __init__(self, **kwargs) -> None:
        self.file = None
        self.validation_results = kwargs.get("validation_results", [])
        self._installer_package = kwargs.get("installer_package", None)
        self._function_app_name = kwargs.get("function_app_name", None)

    @property
    def function_app_name(self):
        return self._function_app_name

    @property
    def installer_package(self):
        return self._installer_package


class CreateApplicationPackageOptions:
    """
    Options for creating an application package.

    Args:
        installer_version (InstallerVersion | str): The version of the installer to use.
        vmi_reference (bool, optional): Whether to use a VMI reference of the published/released reference. Defaults to False.
        vmi_reference_id (str, optional): The ID of the VMI reference to use to override the published reference.
    """
    def __init__(self, installer_version: InstallerVersion | str, vmi_reference: bool = False, vmi_reference_id: str = None) -> None:
        self._use_vmi_reference = vmi_reference

        if isinstance(installer_version, str):
            self.installer_version = InstallerVersion(installer_version)
        else:
            self.installer_version = installer_version
        
        if vmi_reference_id is not None:
            self._vmi_reference_id = vmi_reference_id
            self._use_vmi_reference = True

    @property
    def vmi_reference_id(self):
        """This is the ID of the VMI reference to use to override the published reference."""
        return self._vmi_reference_id
    
    @property
    def use_vmi_reference(self):
        return self._use_vmi_reference


class ApplicationPackageInfo(Model):
    def __init__(self, main_template: str | Path, create_ui_definition: str | CreateUiDefinition, name="", description=""):
        super().__init__()
        self.create_ui_definition = create_ui_definition

        # not to be confused with the main template of the application package.
        # this is the main template for the app the installer will deploy
        self.main_template = main_template
        if isinstance(create_ui_definition, str) or isinstance(create_ui_definition, Path):
            self.create_ui_definition = CreateUiDefinition.from_file(create_ui_definition)
        else:
            self.create_ui_definition = create_ui_definition

        self.manifest = ManifestInfo(main_template=main_template)
        self.manifest.offer.name = name
        self.manifest.offer.description = description

    @property
    def name(self):
        return self.manifest.offer.name
    
    @property
    def description(self):
        return self.manifest.offer.description
    
    @property
    def template_parameters(self):
        return self.manifest.get_parameters()

    def validate(self):
        template_parameters = self.manifest.get_parameters()
        validation_results = self.manifest.validate()
        validation_results += self.create_ui_definition.validate(template_parameters)

        return validation_results

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
    - function.zip
    - installer.zip
        - manifest.json
        - main.ts (the installer's main template)
        - modules
        - <modules>
    """
    file_name = "app.zip"

    def __init__(self, resourcesProvider: InstallerResourcesProvider) -> None:
        self._resourcesProvider = resourcesProvider

    def create(self, info: ApplicationPackageInfo, options: CreateApplicationPackageOptions, out_dir=None) -> CreateApplicationPackageResult:
        """
        Creates an application package based on the current manifest and UI definition.

        Args:
            info (ApplicationPackageInfo): The application package info.
            out_dir (Optional[str]): The output directory for the application package.
            If not specified, the package will be created in a randomly generated temp directory.

        Returns:
            CreateApplicationPackageResult: A result object containing the validation results and the path to the
            created application package file.
        """
        validation_results = info.validate()

        if len(validation_results) > 0:
            return CreateApplicationPackageResult(validation_results=validation_results)

        installer_package = create_installer_package(info.manifest)

        self._finalize_main_template(info, installer_package, options)
        self._finalize_view_definition(info, options)

        result = CreateApplicationPackageResult(
            installer_package=installer_package, function_app_name=self.main_template.function_app_name
        )
        result.file = self._zip(info, installer_package, options, out_dir)

        if result.file is None or not result.file.exists():
            result.validation_results.append(Exception("Failed to create application package"))
            return result

        return result

    def _finalize_view_definition(self, info: ApplicationPackageInfo, options: CreateApplicationPackageOptions):
        view_definition = self._resourcesProvider.get(options.installer_version).view_definition
        view_definition.add_input("dashboardUrl", self.main_template.dashboard_url)
        view_definition.add_input("offerName", info.manifest.offer.name)
        view_definition.add_input("offerDescription", info.manifest.offer.description)

        self.view_definition = view_definition

    def _finalize_main_template(
        self, info: ApplicationPackageInfo, installer_package: CreateInstallerPackageResult, options: CreateApplicationPackageOptions
    ):
        """the app package's main template"""
        installer_resources = self._resourcesProvider.get(options.installer_version)

        finalizer = MainTemplateFinalizer(installer_resources.main_template)
        self.main_template = finalizer.finalize(
            template_parameters=info.template_parameters,
            installer_resources=installer_resources,
            installer_package=installer_package,
            use_vmi_reference=options.use_vmi_reference,
            vmi_reference_id=options.vmi_reference_id,
        )

    def _zip(self, info, installer_package, options, out_dir=None) -> Path:
        installer_resources = self._resourcesProvider.get(options.installer_version)

        out_dir = out_dir if out_dir is not None else tempfile.mkdtemp()
        file = Path(out_dir).joinpath(self.file_name)

        with ZipFile(file, "w") as zip_file:
            zip_file.write(installer_resources.function_app_package, installer_resources.function_app_package.name)
            zip_file.write(installer_package.path, installer_package.name)
            zip_file.writestr(MAIN_TEMPLATE_FILE_NAME, self.main_template.to_json())
            zip_file.writestr(VIEW_DEFINITION_FILE_NAME, self.view_definition.to_json())
            zip_file.writestr(CREATE_UI_DEFINITION_FILE_NAME, info.create_ui_definition.to_json())
            zip_file.close()

        return file


def new_application_package() -> ApplicationPackage:
    return ApplicationPackage(InstallerResourcesProvider())