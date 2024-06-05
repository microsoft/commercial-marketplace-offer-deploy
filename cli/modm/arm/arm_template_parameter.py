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

    input_type = input_variable.type

    if input_type == TerraformInputVariableType.string.value:
        parameter_type = ArmTemplateParameterType.securestring if input_variable.sensitive else ArmTemplateParameterType.string

    elif input_type == TerraformInputVariableType.number.value:
        parameter_type = ArmTemplateParameterType.int

    elif input_type == TerraformInputVariableType.bool.value:
        parameter_type = ArmTemplateParameterType.bool

    elif input_type in {TerraformInputVariableType.object.value, TerraformInputVariableType.map.value}:
        if input_variable.sensitive:
            parameter_type = ArmTemplateParameterType.secureobject
        else:
            parameter_type = ArmTemplateParameterType.object

    elif input_type in {TerraformInputVariableType.list.value, TerraformInputVariableType.set.value}:
        parameter_type = ArmTemplateParameterType.array

    else:
        raise ValueError(f"Unsupported Terraform input variable type {input_type}")

    return ArmTemplateParameter(input_variable.name, parameter_type)

