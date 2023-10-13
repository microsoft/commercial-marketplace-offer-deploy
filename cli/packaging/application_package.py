from pathlib import Path
import tempfile
from zipfile import ZipFile
from packaging import ManifestInfo, InstallerPackage
from packaging.azure.create_ui_definition import CreateUiDefinition
from packaging.installer_package import create_installer_package
from packaging.azure import ArmTemplate
import packaging.manifest as manifest
from importlib.resources import files, as_file


MAIN_TEMPLATE_FILE_NAME = "mainTemplate.json"
CREATE_UI_DEFINITION_FILE_NAME = "createUiDefinition.json"
VIEW_DEFINITION_FILE_NAME = "viewDefinition.json"


class ApplicationPackage:
    file_name = "app.zip"

    """
    Represents the app package, e.g. the app.zip
    The installer package (installer.pkg) will reside directly in the app.zip next to
    the installer's mainTemplate.json and createUiDefinition.json, respectively

    """

    def __init__(self, main_template: ArmTemplate, create_ui_definition: CreateUiDefinition, name="", description="") -> None:
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
        - installer.pkg
            - manifest.json
            - main.ts (the installer's main template)
            - modules
            - <modules>
        """

        self.manifest = ManifestInfo(main_template=main_template)
        self.manifest.offer.name = name
        self.manifest.offer.description = description

        self.main_template = self._get_main_template()
        self.create_ui_definition = create_ui_definition

    def create(self):
        template_parameters = self.manifest.get_parameters()

        validation_results = self.manifest.validate()
        validation_results += self.create_ui_definition.validate(template_parameters)

        if len(validation_results) > 0:
            # TODO: create an aggregate exception type that can hold multiple exceptions
            raise ValueError(f"Invalid application package: {validation_results}")
        
        self._finalize_main_template(template_parameters)
        self._zip()

    def _finalize_main_template(self, template_parameters):
        """
        Updates the (installer's) main template with the parameters from the app's main template.
        This results in a flow of: createUiDefinition.json/parameters/outputs --> mainTemplate.json/parameters

        Explanation: 
            This allows the parameters to be passed to the mainTemplate.json/variables/userData 
            so MODM can bootstrap the application with it's parameters when it performs the deployment
        """
        
        self.main_template.insert_parameters(template_parameters)
        for parameter in template_parameters:
            self.main_template.document["variables"]["userData"][parameter.name] = f"[parameters('{parameter.name}')]"

    def _zip(self):
        installer_package = create_installer_package(self.manifest)
        file = Path(tempfile.mkdtemp()).joinpath(self.file_name)

        with ZipFile(file, "w") as zip_file:
            zip_file.write(installer_package, installer_package.name)
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
