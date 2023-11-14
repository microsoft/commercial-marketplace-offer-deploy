from enum import Enum
from pathlib import Path


class SolutionTemplateType(Enum):
    """
    The type of template for the solution template (that modm will use to perform the deployment)

    Attributes:
      bicep (str): supported deployment type of bicep
      arm (str): supported deployment type of arm
      terraform (str): supported deployment type of terraform
    """

    bicep = "bicep"
    arm = "arm"
    terraform = "terraform"

    @staticmethod
    def is_arm(template_file: Path):
        """
        Check if the template file type is arm.

        Returns:
          bool: True if the deployment type is arm, False otherwise
        """
        extension = template_file.suffix.lstrip(".")
        return extension == SolutionTemplateType.arm.value

    @staticmethod
    def is_bicep(template_file: Path):
        """
        Check if the template file type is bicep.

        Returns:
          bool: True if the deployment type is bicep, False otherwise
        """
        extension = template_file.suffix.lstrip(".")
        return extension == SolutionTemplateType.bicep.value

    @staticmethod
    def is_terraform(template_file: Path):
        """
        Check if the template file type is terraform.

        Returns:
          bool: True if the deployment type is terraform, False otherwise
        """
        extension = template_file.suffix.lstrip(".")
        return extension == "tf"

    @staticmethod
    def from_template(solution_template: Path):
        """
        Get the deployment type from the template file.

        Returns:
          TemplateType: The deployment type
        """
        if SolutionTemplateType.is_arm(solution_template):
            return SolutionTemplateType.arm
        elif SolutionTemplateType.is_bicep(solution_template):
            return SolutionTemplateType.bicep
        elif SolutionTemplateType.is_terraform(solution_template):
            return SolutionTemplateType.terraform
        else:
            raise ValueError(f"Unsupported deployment type for template {solution_template}")
