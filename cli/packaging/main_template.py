

from packaging.azure import ArmTemplate
from packaging.azure.function_app import create_function_app_name


class MainTemplate(ArmTemplate):
    """
    This is the main template of the app.zip that will be used to deploy MODM;
    not to be confused with the "main template" for the application itself which will
    reside in the installer package placed into the app.zip
    """
    function_app_name_prefix = "modmfunc"
    function_app_name_variable = "functionAppName"
    vmi_reference_id_variable = "vmiReferenceId"

    def __init__(self, document, vmi_reference_id = None) -> None:
        super().__init__(document)
        self.vmi_reference_id = vmi_reference_id
        self.function_app_name = create_function_app_name(self.function_app_name_prefix)

    @property
    def vmi_reference_id(self):
        return self.document["variables"][self.vmi_reference_id_variable]
    
    @property.setter
    def vmi_reference_id(self, value):
        self.document["variables"][self.vmi_reference_id_variable] = value
    
    @property
    def dashboard_url(self):
        return f'https://{self.function_app_name}.azurewebsites.net/dashboard'

    @property
    def function_app_name(self):
        """The function app name used to create a FunctionApp which will drive the dashboard"""
        return self.document["variables"][self.function_app_name_variable]

    @property.setter
    def function_app_name(self, value):
        self.document["variables"][self.function_app_name_variable] = value
    

def from_file(file_path: str, vmi_reference_id: str):
    instance = MainTemplate.from_file(file_path)
    instance.vmi_reference_id = vmi_reference_id
    return instance
    