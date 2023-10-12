

from enum import Enum


class AzureArmParameterTypes(Enum):
  string = "string"
  object = "object"

  def convert_terraform_type(terraform_variable_type):
    if terraform_variable_type == "string":
      return AzureArmParameterTypes.string.value
    if terraform_variable_type == "object" or terraform_variable_type == "map":
      return AzureArmParameterTypes.object.value
    if terraform_variable_type == "bool":
      return 

class InputParameter