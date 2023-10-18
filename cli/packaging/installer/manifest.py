import json
from pathlib import Path
from msrest.serialization import Model
import copy
import os
import jsonschema
from pathlib import Path
from packaging.terraform import TerraformFile
from packaging.azure import ArmTemplateParameter, from_terraform_input_variable
from .deployment_type import DeploymentType


class ManifestInfo(Model):
    _attribute_map = {
        "main_template": {"key": "mainTemplate", "type": "str"},
        "deployment_type": {"key": "deploymentType", "type": "str"},
        "offer": {"key": "offer", "type": "OfferInfo"},
    }

    def __init__(self, main_template: Path, **kwargs):
        super().__init__(**kwargs)

        self.main_template = main_template
        self.deployment_type = kwargs.get("deployment_type", "")
        self.offer = OfferInfo()

        if self.deployment_type == "":
            template = Path(self.main_template).name
            if template.endswith(".tf"):
                self.deployment_type = DeploymentType.terraform
            elif template.endswith(".bicep"):
                self.deployment_type = DeploymentType.bicep

    def to_json(self):
        return json.dumps(self.serialize(), indent=2)

    def get_parameters(self) -> list[ArmTemplateParameter]:
        """
        Returns the parameters of the app's main template as a list of ArmTemplateParameter
        """
        if self.deployment_type == DeploymentType.terraform:
            terraform_file = TerraformFile(self.main_template)
            input_variables = terraform_file.parse_variables()
            parameters = list(map(from_terraform_input_variable, input_variables))

            return parameters
        else:
            # TODO: support bicep parameters. We'll need to convert the bicep template to ARM template first, then
            # extract the parameters from the ARM template, directly
            raise ValueError(f"Unsupported deployment type {self.deployment_type}")

    def validate(self):
        validation_results = super().validate()

        if validation_results is None:
            validation_results = []

        # validate the app's main template exists and is matching the deployment type
        main_template_file = Path(self.main_template)

        if not main_template_file.exists():
            validation_results.append(FileNotFoundError(f"Could not find main template file at {self.main_template}"))

        if not main_template_file.is_file():
            validation_results.append(FileNotFoundError(f"Main template file {self.main_template} is not a file"))

        if self.deployment_type == DeploymentType.terraform and main_template_file.suffix != ".tf":
            validation_results.append(ValueError(f"Main template file {self.main_template} must have a .tf extension"))

        return validation_results


class OfferInfo(Model):
    _attribute_map = {"name": {"key": "name", "type": "str"}, "description": {"key": "description", "type": "str"}}

    def __init__(self, **kwargs):
        super().__init__(**kwargs)

        self.name = kwargs.get("name", "")
        self.description = kwargs.get("description", "")


class ManifestFile:
    file_name = "manifest.json"

    @staticmethod
    def write(dest_path, manifest: ManifestInfo):
        manifest_copy = copy.deepcopy(manifest)
        manifest_copy.main_template = os.path.basename(manifest.main_template)

        json = manifest_copy.to_json()
        file_path = Path(os.path.join(dest_path, ManifestFile.file_name)).resolve()

        with open(file_path, "w") as f:
            f.write(json)

    def validate(self, schema):
        try:
            jsonschema.validate(self.read(), schema)
        except jsonschema.exceptions.ValidationError as err:
            raise err


def write_manifest(dest_path, manifest: ManifestInfo):
    ManifestFile.write(dest_path, manifest)
