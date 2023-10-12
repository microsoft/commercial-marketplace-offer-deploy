from pathlib import Path
from packaging import ManifestInfo, InstallerPackage
from packaging.azure import ArmTemplate
import packaging.manifest as manifest
from importlib.resources import files, as_file

class ApplicationPackage:
  """
  Represents the app package, e.g. the app.zip
  The installer package (installer.pkg) will reside directly in the app.zip next to 
  the mainTemplate.json and createUiDefinition.json, respectively

  """

  def __init__(self, main_template, create_ui_definition, name='', description='') -> None:
    """
    Args:
      main_template (str): The path to the application's main template (the one that).
      name (str, optional): The name of the application. Defaults to ''.
      description (str, optional): The description of the application. Defaults to ''.

    Example Output (structure):
    - app.zip
      - mainTemplate.json
      - createUiDefinition.json
      - viewDefinition.json
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

  

  def create_ui_definition(self, main_template):
    """
    The createUiDefinition.json that will be used to create the Azure Portal experience.
    The outputs of this file MUST match the input parameters of the application_main_template
    """
    self.manifest.main_template = main_template
    return self

  def create(self):
    installer_package = InstallerPackage(self.manifest)
    package_file, package_dir = installer_package.create()
    

  def _get_main_template(self):
    resource_files = files("resources.templates")
    
    # this is the main template of the app.zip that will be used to deploy MODM. 
    # Not to be confused with the main template for the application itself which will
    # reside in the installer package contained in app.zip
    with as_file(resource_files.joinpath("mainTemplate.json")) as resource_file:
      template = ArmTemplate.from_file(resource_file)
    return template