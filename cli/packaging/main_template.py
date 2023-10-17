

from packaging.azure import ArmTemplate


class MainTemplate(ArmTemplate):
    """
    This is the main template of the app.zip that will be used to deploy MODM;
    not to be confused with the "main template" for the application itself which will
    reside in the installer package placed into the app.zip
    """
    def __init__(self, document, vmi_reference_id = None) -> None:
        super().__init__(document)
        self.vmi_reference_id = vmi_reference_id
        self.function_app_name = None

    @property
    def vmi_reference_id(self):
        return self.document["variables"]["functionAppName"]
    @property
    def function_app_name(self):
        return self.document["variables"]["functionAppName"]
    

def from_file(file_path: str, vmi_reference_id: str):
    instance = MainTemplate.from_file(file_path)
    instance.vmi_reference_id = vmi_reference_id
    return instance
    