import json
import os
from modm.arm.arm_template import ArmTemplate
from modm.arm.arm_template_parameter import ArmTemplateParameter
from modm.installer.reserved_template_parameter import ReservedTemplateParameter


class UserData:
    installer_package_key = "installerPackage"

    def __init__(self, document: dict) -> None:
        self.document = document

    def set_installer_package_hash(self, value):
        if self.installer_package_key not in self.document:
            self.document[self.installer_package_key] = {}
            self.document[self.installer_package_key][
                "uri"
            ] = "[concat(variables('artifactsContainerLocation'), '/', 'installer.zip')]"

        installer_package = self.document[self.installer_package_key]
        installer_package["hash"] = value

    def insert_parameters(self, parameters: list[ArmTemplateParameter]):
        user_data_parameters = self.document["parameters"]
        for parameter in parameters:
            user_data_parameters[parameter.name] = f"[parameters('{parameter.name}')]"


class MainTemplate(ArmTemplate):
    """
    This is the main template of the app.zip that will be used to deploy MODM;
    not to be confused with the solution template for the application which will
    reside in the installer package (installer.zip) included in the app.zip
    """

    vmi_reference_id_variable = "vmiReferenceId"
    vm_storage_profile_image_reference = "imageReference"

    def __init__(self, document) -> None:
        super().__init__(document)
        self._vm_offer = None
        self._user_data = UserData(self.document["variables"]["userData"])

    def set_parameters(self, parameters: list[ArmTemplateParameter]):
        reserved_parameters = ReservedTemplateParameter.all()
        for parameter in parameters:
            if parameter.name in reserved_parameters:
                parameters.remove(parameter)

        super().set_parameters(parameters)
        self._user_data.insert_parameters(parameters)

    @property
    def vmi_reference_id(self):
        return self.document["variables"][self.vmi_reference_id_variable]

    @vmi_reference_id.setter
    def vmi_reference_id(self, value):
        self.document["variables"][self.vmi_reference_id_variable] = value

        resources = self.document["resources"]
        for resource in resources:
            if resource["type"] == "Microsoft.Compute/virtualMachines":
                resource["properties"]["storageProfile"]["imageReference"] = {"id": "[variables('vmiReferenceId')]"}
                break

    @property
    def vm_offer(self):
        return self._vm_offer

    @vm_offer.setter
    def vm_offer(self, value: dict):
        self._vm_offer = value

        resources = self.document["resources"]
        for resource in resources:
            if resource["type"] == "Microsoft.Compute/virtualMachines":
                resource["plan"] = value["plan"]
                resource["properties"]["storageProfile"]["imageReference"] = value["imageReference"]
                break

        if self.vmi_reference_id_variable in self.document["variables"]:
            del self.document["variables"][self.vmi_reference_id_variable]

    @property
    def user_data(self) -> UserData:
        return self._user_data

    @staticmethod
    def from_file(file_path):
        """
        Load an ARM template from a file.

        Args:
          file_path (str): The path to the ARM template file.

        Returns:
          ArmTemplate: An instance of the ArmTemplate class representing the loaded ARM template.
        """
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Could not find ARM template file at {file_path}")

        with open(file_path, "r") as f:
            document = json.load(f)
            return MainTemplate(document)


def from_file(file_path: str) -> MainTemplate:
    instance = MainTemplate.from_file(file_path)
    return instance

