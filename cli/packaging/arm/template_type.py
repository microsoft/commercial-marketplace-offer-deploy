from enum import Enum
from pathlib import Path


class TemplateType(Enum):
    """
    The type of template

    Attributes:
      bicep (str): supported deployment type of bicep
      arm (str): supported deployment type of arm
    """

    bicep = "bicep"
    arm = "arm"

    @staticmethod
    def is_arm(template_file: Path):
        """
        Check if the template file type is arm.

        Returns:
          bool: True if the deployment type is arm, False otherwise
        """
        extension = template_file.suffix.lstrip(".")
        return extension == TemplateType.arm.value

    @staticmethod
    def is_bicep(template_file: Path):
        """
        Check if the template file type is bicep.

        Returns:
          bool: True if the deployment type is bicep, False otherwise
        """
        extension = template_file.suffix.lstrip(".")
        return extension == TemplateType.bicep.value