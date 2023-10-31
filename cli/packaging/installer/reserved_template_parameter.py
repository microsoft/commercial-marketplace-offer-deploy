
from enum import Enum

class ReservedTemplateParameter(str, Enum):
    """
    The reserved template parameters for the installer. These are the parameters that are provided by the installer.
    They should NOT be provided by the user in the createUiDefinition.json file.
    """
    resource_group_name = "resourceGroupName"
    
    @classmethod
    def all(cls):
        return [member.value for member in cls]


def is_reserved(value: str) -> bool:
    """
    Check if a given value is a reserved template parameter.

    Args:
        value (str): The value to check.

    Returns:
        bool: True if the value is a reserved template parameter, False otherwise.
    """
    return value in ReservedTemplateParameter.all()
