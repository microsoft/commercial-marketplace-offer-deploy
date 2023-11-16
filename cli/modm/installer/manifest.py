import json
from pathlib import Path
import shutil
import tempfile
from msrest.serialization import Model
import copy
import os
import jsonschema
from pathlib import Path
from modm.arm.bicep_template_compiler import BicepTemplateCompiler
from modm.terraform import TerraformFile
from modm.arm import ArmTemplateParameter, from_terraform_input_variable
from modm.arm.arm_template import ArmTemplate
from modm.installer.solution_template_type import SolutionTemplateType
from .deployment_type import DeploymentType


class ManifestInfo(Model):
    """
    This is the manifest.json file that will be included in the installer.zip
    """

    _attribute_map = {
        "solution_template": {"key": "mainTemplate", "type": "str"},
        "deployment_type": {"key": "deploymentType", "type": "str"},
        "offer": {"key": "offer", "type": "OfferProperties"},
    }

    def __init__(self, solution_template: Path, **kwargs):
        super().__init__(**kwargs)

        self.solution_template = Path(solution_template)
        self._template_type = SolutionTemplateType.from_template(self.solution_template)

        self._bicep_templates_dir = None
        self._compile_bicep_template()

        self.offer = OfferProperties()

        if self._template_type == SolutionTemplateType.terraform:
            self.deployment_type = DeploymentType.terraform
        else:
            self.deployment_type = DeploymentType.arm

    @property
    def template_type(self) -> SolutionTemplateType:
        return self._template_type
    
    @property 
    def has_bicep_source(self) -> Path:
        return self._bicep_templates_dir is not None and self._bicep_templates_dir.exists()

    @property 
    def bicep_templates_dir(self) -> Path:
        return self._bicep_templates_dir

    def to_json(self):
        return json.dumps(self.serialize(), indent=2)

    def get_parameters(self) -> list[ArmTemplateParameter]:
        """
        Returns the parameters of the app's main template as a list of ArmTemplateParameter
        """
        if self.template_type == SolutionTemplateType.terraform:
            terraform_file = TerraformFile(self.solution_template)
            input_variables = terraform_file.parse_variables()
            parameters = list(map(from_terraform_input_variable, input_variables))

            return parameters
        elif self.template_type == SolutionTemplateType.arm:
            arm_template = ArmTemplate.from_file(self.solution_template)
            return arm_template.get_parameters()
        else:
            raise ValueError(f"Unsupported template type {self.template_type}")

    def validate(self):
        validation_results = super().validate()

        if validation_results is None:
            validation_results = []

        # validate the app's main template exists and is matching the deployment type
        main_template_file = Path(self.solution_template)

        if not main_template_file.exists():
            validation_results.append(FileNotFoundError(f"Could not find main template file at {self.solution_template}"))

        if not main_template_file.is_file():
            validation_results.append(FileNotFoundError(f"Main template file {self.solution_template} is not a file"))

        if self.deployment_type == DeploymentType.terraform and main_template_file.suffix != ".tf":
            validation_results.append(ValueError(f"Main template file {self.solution_template} must have a .tf extension"))

        return validation_results

    def _compile_bicep_template(self):
        if self.template_type == SolutionTemplateType.bicep:
            src_templates_dir = self.solution_template.parent
            temp_dir = Path(tempfile.mkdtemp())

            self._bicep_templates_dir = temp_dir / ".bicep"
            self._bicep_templates_dir.mkdir()
            shutil.copytree(str(src_templates_dir), str(self._bicep_templates_dir), dirs_exist_ok=True)

            compiler = BicepTemplateCompiler(self.solution_template)
            arm_template_file = compiler.compile(temp_dir)

            # update the solution template we're pointing to since it's now the compiled arm template
            self.solution_template = arm_template_file
            self._template_type = SolutionTemplateType.arm

    def dispose(self):
        if self.has_bicep_source:
            shutil.rmtree(self._bicep_templates_dir)

class OfferProperties(Model):
    """
    This is information about the publisher's offer NOT the app installer's offer. 
    """
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
        manifest_copy.solution_template = os.path.basename(manifest.solution_template)

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
