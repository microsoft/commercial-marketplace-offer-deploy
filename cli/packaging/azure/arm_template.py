import os
import json

class ArmTemplate:
  def __init__(self, document):
    self.document = document
    self.name = "mainTemplate.json"

  def write(self, file_path):
    # Create directory if it doesn't exist
    os.makedirs(os.path.dirname(file_path), exist_ok=True)

    # Write formatted JSON to file
    with open(file_path, 'w') as f:
      json.dump(self.document, f, indent=4)

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
      raise FileNotFoundError(f'Could not find ARM template file at {file_path}')

    with open(file_path, 'r') as f:
      document = json.load(f)
      return ArmTemplate(document)
