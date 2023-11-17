from enum import Enum


class DeploymentType(Enum):
    """
    The type of deployment.

    Attributes:
      terraform (str): supported deployment type of terraform
      arm (str): supported deployment type of arm
    """

    terraform = "terraform"
    arm = "arm"
