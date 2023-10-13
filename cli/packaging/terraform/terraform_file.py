import hcl2
from .input_variable import TerraformInputVariable


class TerraformFile(object):
    """
    A class representing a Terraform file.

    Attributes:
      file_path (str): The path to the Terraform file.
    """

    def __init__(self, file_path):
        self.file_path = file_path

    def parse_variables(self) -> list[TerraformInputVariable]:
        with open(self.file_path, "r") as file:
            dict = hcl2.load(file)

            if "variable" in dict:
                variables = dict["variable"]
                return list(map(TerraformInputVariable, variables))
            # no variables were present in the template
            return [TerraformInputVariable]
