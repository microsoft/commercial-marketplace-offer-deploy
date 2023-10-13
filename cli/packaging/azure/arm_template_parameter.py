from enum import Enum
from packaging.terraform.input_variable import TerraformInputVariable
from packaging.terraform.variable_types import TerraformInputVariableType


class ArmTemplateParameter:
    def __init__(self, name, type):
        self.name = name
        self.type = type


class ArmTemplateParameterType(Enum):
    string = "string"
    securestring = "secureString"
    int = "int"
    bool = "bool"
    object = "object"
    secureobject = "secureObject"
    array = "array"


def from_terraform_input_variable(input_variable: TerraformInputVariable) -> ArmTemplateParameter:
    parameter_type = None

    if input_variable.type == TerraformInputVariableType.string.value:
        if input_variable.sensitive:
            parameter_type = ArmTemplateParameterType.securestring.value
        else:
            parameter_type = ArmTemplateParameterType.string
    
    if input_variable.type == TerraformInputVariableType.number.value:
        parameter_type = ArmTemplateParameterType.int

    if input_variable.type == TerraformInputVariableType.bool.value:
        parameter_type = ArmTemplateParameterType.bool

    if input_variable.type == TerraformInputVariableType.object or input_variable.type == TerraformInputVariableType.map.value:
        if input_variable.sensitive:
            parameter_type = ArmTemplateParameterType.secureobject
        else:
            parameter_type = ArmTemplateParameterType.object
    
    if input_variable.type == TerraformInputVariableType.list.value or input_variable.type == TerraformInputVariableType.set.value:
        parameter_type = ArmTemplateParameterType.array

    if parameter_type is None:
        raise ValueError(f"Unsupported Terraform input variable type {input_variable.type}")
    
    return ArmTemplateParameter(input_variable.name, parameter_type)
