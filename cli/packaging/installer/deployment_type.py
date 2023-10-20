from enum import Enum


class DeploymentType(Enum):
    """
    An enumeration of Terraform variable types.

    Attributes:
      terraform (str): supported deployment type of terraform
      bicep (str): supported deployment type of bicep
    """

    terraform = "terraform"
    bicep = "bicep"
