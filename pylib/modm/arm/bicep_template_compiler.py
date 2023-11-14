import os
import re
import subprocess
from pathlib import Path
from knack.log import get_logger

_bicep_diagnostic_warning_pattern = r"^([^\s].*)\((\d+)(?:,\d+|,\d+,\d+)?\)\s+:\s+(Warning)\s+([a-zA-Z-\d]+):\s*(.*?)\s+\[(.*?)\]$"  # pylint: disable=line-too-long

class BicepTemplateCompiler:
    def __init__(self, file_path: Path):
        self.file_path = file_path
        self._logger = get_logger(__name__)

    @property
    def file_name(self) -> str:
        return self.file_path.stem

    def compile(self, out_dir: Path) -> Path:
        compiled_file_path = out_dir / (self.file_name + ".json")

        process = subprocess.run(
            ["az", "bicep", "build", "--file", str(self.file_path), "--outdir", str(out_dir)],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=None)

        try:
            process.check_returncode()
            command_warnings = process.stderr.decode("utf-8")
            if command_warnings:
                self._logger.warning(command_warnings)

            self._logger.debug(process.stdout.decode("utf-8"))
            return Path(compiled_file_path)

        except subprocess.CalledProcessError:
            stderr_output = process.stderr.decode("utf-8")
            errors = []

            for line in stderr_output.splitlines():
                if re.match(_bicep_diagnostic_warning_pattern, line):
                    self._logger.warning(line)
                else:
                    errors.append(line)

            error_msg = os.linesep.join(errors)
            raise Exception(error_msg)
