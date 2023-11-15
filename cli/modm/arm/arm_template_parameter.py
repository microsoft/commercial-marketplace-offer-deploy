from enum import Enum
from modm.terraform.input_variable import TerraformInputVariable
from modm.terraform.variable_types import TerraformInputVariableType


class ArmTemplateParameter:
    def __init__(self, name: str, type):
        self.name = name
        self.type = type
        self.default_value = None

    def value(self):
        """Returns the ARM template parameter value as a dict"""
        if isinstance(self.type, Enum):
            type_value = self.type.value
        else:
            type_value = self.type
        value = {
            "type": type_value
        }
        if self.default_value is not None:
            value["defaultValue"] = self.default_value

        return value


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
            parameter_type = ArmTemplateParameterType.securestring
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
