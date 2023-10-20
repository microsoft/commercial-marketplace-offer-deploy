from pathlib import Path
import tempfile
from zipfile import ZipFile
from packaging.application_delivery import ApplicationDelivery
from packaging.installer import ManifestInfo, CreateInstallerPackageResult, create_installer_package
import packaging.azure as azure
from packaging.azure import CreateUiDefinition
from packaging.function_app_package import FunctionAppPackage
from importlib.resources import files, as_file
from msrest.serialization import Model

from packaging.installer.resources import InstallerResources, InstallerResourcesProvider
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
    dashboard_url_user_data_variable = "dashboardUrl"
    function_app_name_variable = "functionAppName"
    vmi_reference_id_variable = "vmiReferenceId"
    default_function_app_name_prefix = "modmfunc"

    def __init__(
        self, installer_version: InstallerVersion | str, delivery_type: ApplicationDelivery = ApplicationDelivery.MARKETPLACE
    ) -> None:
        if isinstance(installer_version, str):
            self.installer_version = InstallerVersion(installer_version)
        else:
            self.installer_version = installer_version

        self.delivery_type = delivery_type
        self.function_app_name = azure.create_function_app_name(self.default_function_app_name_prefix)

    @property
    def dashboard_url(self):
        return f"https://{self.function_app_name}.azurewebsites.net/dashboard"

    def get_installer_resources(self) -> InstallerResources:
        return InstallerResourcesProvider().get(self.installer_version.name)


class ApplicationPackage:
    file_name = "app.zip"

    """
    Represents the app package, e.g. the app.zip
    The installer package (installer.zip) will reside directly in the app.zip next to
    the installer's mainTemplate.json and createUiDefinition.json, respectively

    """

    def __init__(self, main_template: str, create_ui_definition: str | CreateUiDefinition, name="", description="") -> None:
        """
        Args:
        main_template (str): The path to the application's main template.
        name (str, optional): The name of the application. Defaults to ''.
        description (str, optional): The description of the application. Defaults to ''.

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

        self.manifest = ManifestInfo(main_template=main_template)
        self.manifest.offer.name = name
        self.manifest.offer.description = description

        self.main_template = None
        self.view_definition = None

        if isinstance(create_ui_definition, str) or isinstance(create_ui_definition, Path):
            self.create_ui_definition = CreateUiDefinition.from_file(create_ui_definition)
        else:
            self.create_ui_definition = create_ui_definition

    def create(self, options: CreateApplicationPackageOptions, out_dir=None) -> CreateApplicationPackageResult:
        """
        Creates an application package based on the current manifest and UI definition.

        Args:
            out_dir (Optional[str]): The output directory for the application package.
            If not specified, the package will be created in a randomly generated temp directory.

        Returns:
            CreateApplicationPackageResult: A result object containing the validation results and the path to the
            created application package file.
        """

        template_parameters = self.manifest.get_parameters()

        validation_results = self.manifest.validate()
        validation_results += self.create_ui_definition.validate(template_parameters)

        if len(validation_results) > 0:
            return CreateApplicationPackageResult(validation_results=validation_results)

        installer_package = create_installer_package(self.manifest)

        self._finalize_main_template(template_parameters, installer_package, options)
        self._finalize_view_definition(options)

        result = CreateApplicationPackageResult(installer_package=installer_package, function_app_name=options.function_app_name)
        result.file = self._zip(installer_package, options, out_dir)

        if result.file is None or not result.file.exists():
            result.validation_results.append(Exception("Failed to create application package"))
            return result

        return result

    def _finalize_view_definition(self, options: CreateApplicationPackageOptions):
        view_definition = options.get_installer_resources().view_definition

        view_definition.add_input("dashboardUrl", options.dashboard_url)
        view_definition.add_input("offerName", self.manifest.offer.name)
        view_definition.add_input("offerDescription", self.manifest.offer.description)

        self.view_definition = view_definition

    def _finalize_main_template(
        self, template_parameters, installer_package: CreateInstallerPackageResult, options: CreateApplicationPackageOptions
    ):
        """
        Updates the (installer's) main template with the parameters from the app's main template.
        This results in a flow of: createUiDefinition.json/parameters/outputs --> mainTemplate.json/parameters

        Explanation:
            This allows the parameters to be passed to the mainTemplate.json/variables/userData
            so MODM can bootstrap the application with it's parameters when it performs the deployment
        """
        installer_resources = options.get_installer_resources()

        main_template = installer_resources.main_template
        main_template.insert_parameters(template_parameters)

        main_template.function_app_name = options.function_app_name
        main_template.user_data.dashboard_url = options.dashboard_url
        main_template.user_data.set_installer_package_hash(installer_package.hash)

        # depending on the delivery type, the template will either use the VM offer or the VMI reference ID
        if options.delivery_type == ApplicationDelivery.MARKETPLACE:
            main_template.vm_offer = installer_resources.vm_offer
        else:
            main_template.vmi_reference_id = installer_resources.vmi_reference_id

        self.main_template = main_template

    def _zip(self, installer_package, options: CreateApplicationPackageOptions, out_dir=None) -> Path:
        out_dir = out_dir if out_dir is not None else tempfile.mkdtemp()
        file = Path(out_dir).joinpath(self.file_name)

        with ZipFile(file, "w") as zip_file:
            zip_file.write(options.function_app_package.path, options.function_app_package.name)
            zip_file.write(installer_package.path, installer_package.name)
            zip_file.writestr(MAIN_TEMPLATE_FILE_NAME, self.main_template.to_json())
            zip_file.writestr(VIEW_DEFINITION_FILE_NAME, self.view_definition.to_json())
            zip_file.writestr(CREATE_UI_DEFINITION_FILE_NAME, self.create_ui_definition.to_json())
            zip_file.close()

        return file
