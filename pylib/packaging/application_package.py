from packaging import ManifestInfo, InstallerPackage
import packaging.manifest as manifest


class ApplicationPackage:
  """
  Represents the app package, e.g. the app.zip
  The installer package (installer.pkg) will reside directly in the app.zip next to 
  the mainTemplate.json and createUiDefinition.json, respectively

  """

  def __init__(self, deployment_type, application_main_template) -> None:
    """
    The application's main template (which is NOT the mainTemplate.json that will be in the root of the app.zip)
    """
    self.manifest = ManifestInfo(deployment_type=deployment_type, main_template=application_main_template)
    

  def main_template(self, main_template):
    """
    The mainTemplate.json that installs modm
    """
    self.manifest.main_template = main_template
    return self
  
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
    

