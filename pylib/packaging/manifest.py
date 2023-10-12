
import json
from pathlib import Path
from msrest.serialization import Model

from packaging.deployment_type import DeploymentType

class ManifestInfo(Model):
  _attribute_map = {
    "main_template": {"key": "mainTemplate", "type": "str"},
    "deployment_type": {"key": "deploymentType", "type": "str"},
    "offer": {"key": "offer", "type": "OfferInfo"},
  }

  def __init__(self, **kwargs):
    super().__init__(**kwargs)

    self.main_template = kwargs.get('main_template', '')
    self.deployment_type = kwargs.get('deployment_type', '')
    self.offer = OfferInfo()
  
  def to_json(self):
    return json.dumps(self.serialize(), indent=2)
  
  def validate(self):
    validation_results = super().validate()
    
    if validation_results is None:
      validation_results = []

    main_template_file = Path(self.main_template)
    
    if not main_template_file.exists():
      validation_results.append(FileNotFoundError(f'Could not find main template file at {self.main_template}'))
    
    if not main_template_file.is_file():
      validation_results.append(FileNotFoundError(f'Main template file {self.main_template} is not a file'))
    
    if self.deployment_type == DeploymentType.terraform and main_template_file.suffix != '.tf':
      validation_results.append(ValueError(f'Main template file {self.main_template} must have a .tf extension'))

    return validation_results

class OfferInfo(Model):
  _attribute_map = {
    "name": {"key": "name", "type": "str"},
    "description": {"key": "description", "type": "str"}
  }

  def __init__(self, **kwargs):
    super().__init__(**kwargs)
    
    self.name = kwargs.get('name', '')
    self.description = kwargs.get('description', '')