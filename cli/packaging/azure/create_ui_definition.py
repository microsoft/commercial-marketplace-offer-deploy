import json
import os
from .arm_template import ArmTemplate


class CreateUiDefinition:
  def __init__(self, document):
    self.document = document

  def validate(self, template_input_parameters: dict):
    validation_results = []
    outputs = self.document['parameters']['outputs']

    if outputs is None:
      validation_results.append(ValueError('The createUiDefinition.json must contain an outputs section'))
      return validation_results

    outputs_keys = set(outputs.keys())
    template_input_parameters_keys = set(template_input_parameters.keys())

    diff = outputs_keys.difference(template_input_parameters_keys)

    if len(diff) > 0:
      validation_results.append(ValueError('The number of output values defined in createUiDefinition.json must match the input parameters count of your template.'))

    return validation_results


  @staticmethod
  def from_file(file_path):
    if not os.path.exists(file_path):
      raise FileNotFoundError(f'Could not find ARM template file at {file_path}')

    with open(file_path, 'r') as f:
      document = json.load(f)
      return CreateUiDefinition(document)
