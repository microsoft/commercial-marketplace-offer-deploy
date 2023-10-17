from pathlib import Path
import tempfile
from zipfile import ZipFile
from packaging import ManifestInfo
from packaging import azure
from packaging.azure.create_ui_definition import CreateUiDefinition
from packaging.installer_package import CreateInstallerPackageResult, create_installer_package
from packaging.azure import ArmTemplate
import packaging.manifest as manifest
from importlib.resources import files, as_file
from msrest.serialization import Model

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
        self.validation_results = kwargs.get('validation_results', [])
        self._installer_package = kwargs.get('installer_package', None)
        self._function_app_name = kwargs.get('function_app_name', None)
    
    @property
    def function_app_name(self):
        return self._function_app_name

    @property
    def installer_package(self):
        return self._installer_package

class CreateApplicationPackageOptions:
    function_app_name_variable = 'functionAppName'
    vmi_reference_id_variable = 'vmiReferenceId'
    default_function_app_name_prefix = "modmfunc"

    def __init__(self, vmi_reference_id) -> None:
        self.vmi_reference_id = vmi_reference_id
        self.function_app_name = azure.create_function_app_name(self.default_function_app_name_prefix)

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
        - installerFunction.zip
        - installer.zip
            - manifest.json
            - main.ts (the installer's main template)
            - modules
            - <modules>
        """

        self.manifest = ManifestInfo(main_template=main_template)
        self.manifest.offer.name = name
        self.manifest.offer.description = description

        self.main_template = self._get_main_template()

        if isinstance(create_ui_definition, str) or isinstance(create_ui_definition, Path):
            self.create_ui_definition = CreateUiDefinition.from_file(create_ui_definition)
        else:
            self.create_ui_definition = create_ui_definition

    def create(self, options: CreateApplicationPackageOptions, out_dir = None) -> CreateApplicationPackageResult:
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
            return CreateApplicationPackageResult(validation_results = validation_results)
        
        installer_package = create_installer_package(self.manifest)

        result = CreateApplicationPackageResult(installer_package = installer_package, function_app_name = options.function_app_name)

        self._finalize_main_template(template_parameters, installer_package, options)
        result.file = self._zip(installer_package, out_dir)

        if result.file is None or not result.file.exists():
            result.validation_results.append(Exception("Failed to create application package"))
            return result
        
        return result

    def _finalize_main_template(self, template_parameters, installer_package: CreateInstallerPackageResult, 
                                options: CreateApplicationPackageOptions):
        """
        Updates the (installer's) main template with the parameters from the app's main template.
        This results in a flow of: createUiDefinition.json/parameters/outputs --> mainTemplate.json/parameters

        Explanation: 
            This allows the parameters to be passed to the mainTemplate.json/variables/userData 
            so MODM can bootstrap the application with it's parameters when it performs the deployment
        """
        
        self.main_template.insert_parameters(template_parameters)

        # variables
        self.main_template.document["variables"][options.function_app_name_variable] = options.function_app_name
        self.main_template.document["variables"][options.vmi_reference_id_variable] = options.vmi_reference_id

        # verify that the userData's installerPackage uri is set
        if "installerPackage" not in self.main_template.document["variables"]["userData"]:
            raise ValueError("failed to identify installerPackage property on userData. Invalid mainTemplate.json")
        
        self.main_template.document["variables"]["userData"]["installerPackage"]['hash'] = installer_package.hash

        # parameters to userData
        self.main_template.document["variables"]["userData"]["parameters"] = {}
        user_data_parameters = self.main_template.document["variables"]["userData"]["parameters"]

        for parameter in template_parameters:
            user_data_parameters[parameter.name] = f"[parameters('{parameter.name}')]"

    def _zip(self, installer_package, out_dir = None) -> Path:
        file = Path(out_dir if out_dir is not None else tempfile.mkdtemp()).joinpath(self.file_name)

        with ZipFile(file, "w") as zip_file:
            zip_file.write(installer_package.path, installer_package.name)
            zip_file.writestr(MAIN_TEMPLATE_FILE_NAME, self.main_template.to_json())
            zip_file.writestr(CREATE_UI_DEFINITION_FILE_NAME, self.create_ui_definition.to_json())
            zip_file.close()
        
        return file

    def _get_main_template(self):
        """
        This is the main template of the app.zip that will be used to deploy MODM;
        not to be confused with the main template for the application itself which will
        reside in the installer package contained in app.zip
        """
        resource_files = files("resources.templates")
        with as_file(resource_files.joinpath("mainTemplate.json")) as resource_file:
            template = ArmTemplate.from_file(resource_file)
            return template
