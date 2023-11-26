import json
from .create_ui_definition import CreateUiDefinition
from modm.installer import ManifestInfo
from msrest.serialization import Model
from pathlib import Path


class ApplicationPackageInfo(Model):
    def __init__(self, solution_template: str | Path, create_ui_definition: str | CreateUiDefinition, name="", description=""):
        """
        Initializes a new instance of the ApplicationPackage class.

        Args:
            solution_template (str | Path): The path to the main template for the app the installer will deploy (NOT modm's mainTemplate.json)
            create_ui_definition (str | CreateUiDefinition): The path to the create UI definition file, or a CreateUiDefinition object.
            name (str, optional): The name of the offer. Defaults to "".
            description (str, optional): The description of the offer. Defaults to "".
        """
        super().__init__()
        self.create_ui_definition = create_ui_definition

        # not to be confused with the main template of the application package.
        # this is the main template for the app the installer will deploy
        self.solution_template = solution_template

        if isinstance(create_ui_definition, str) or isinstance(create_ui_definition, Path):
            self.create_ui_definition = CreateUiDefinition.from_file(create_ui_definition)
        else:
            self.create_ui_definition = create_ui_definition

        self.manifest = ManifestInfo(solution_template=solution_template)
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
    
    @property
    def template_type(self):
        return self.manifest.template_type
    
    def print_manifest_as_json(self):
        manifest_dict = self.manifest.as_dict()
        manifest_json = json.dumps(manifest_dict, indent=2)
        print(manifest_json)

    def validate(self):
        template_parameters = self.manifest.get_parameters()
        validation_results = self.manifest.validate()
        validation_results += self.create_ui_definition.validate(template_parameters)

        return validation_results
    