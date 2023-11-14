from enum import Enum
from pathlib import Path


class DeploymentType(Enum):
    """
    The type of deployment.

    Attributes:
      terraform (str): supported deployment type of terraform
      arm (str): supported deployment type of arm
    """

    terraform = "terraform"
    arm = "arm"

    @staticmethod
    def is_arm(template_file: Path):
        """
        Check if the deployment type is arm.

        Returns:
          bool: True if the deployment type is arm, False otherwise
        """
        extension = template_file.suffix.lstrip(".")
        return extension == DeploymentType.arm.value or extension == "bicep"