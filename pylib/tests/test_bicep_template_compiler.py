import json
import logging
import os
import shutil
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch, call, MagicMock
from modm.arm.bicep_template_compiler import BicepTemplateCompiler

log = logging.getLogger(__name__)

class TestBicepTemplateCompiler(unittest.TestCase):
    def setUp(self):
        self.data_path = os.path.join(os.path.dirname(__file__), 'data')
        self.main_template_file = Path(os.path.join(self.data_path, 'simple_bicep', 'main.bicep'))
        self.compiler = BicepTemplateCompiler(self.main_template_file)

    def test_name(self):
        self.assertEqual(self.compiler.file_name, "main")

    def test_compile_live(self):
        out_dir = Path(tempfile.mkdtemp())
        
        compiled_file_path = self.compiler.compile(out_dir)

        self.assertEqual(compiled_file_path, out_dir / "main.json")
        self.assertTrue(compiled_file_path.exists())

        json_document = json.loads(compiled_file_path.read_text())
        self.assertTrue("parameters" in json_document)
        self.assertTrue("location" in json_document["parameters"])

        shutil.rmtree(out_dir)

    @patch("modm.arm.bicep_template_compiler.subprocess.run")
    def test_compile(self, mock_run):
        mock_process = MagicMock()
        mock_process.returncode = 0
        mock_process.stderr.decode.return_value = ""
        mock_process.stdout.decode.return_value = "Compilation succeeded."
        mock_run.return_value = mock_process

        out_dir = Path("out")
        compiled_file_path = self.compiler.compile(out_dir)

        mock_run.assert_called_once_with(
            ["az", "bicep", "build", "--file", str(self.main_template_file), "--outdir", str(out_dir)],
            stdout=-1,
            stderr=-1,
            env=None
        )

        self.assertEqual(compiled_file_path, out_dir / "main.json")

    @patch("modm.arm.bicep_template_compiler.subprocess.run")
    def test_compile_with_warnings(self, mock_run):
        mock_process = MagicMock()
        mock_process.returncode = 0
        mock_process.stderr.decode.return_value = "Warning: unused variable 'foo'"
        mock_process.stdout.decode.return_value = "Compilation succeeded."
        mock_run.return_value = mock_process

        out_dir = Path("out")
        compiled_file_path = self.compiler.compile(out_dir)

        mock_run.assert_called_once_with(
            ["az", "bicep", "build", "--file", str(self.main_template_file), "--outdir", str(out_dir)],
            stdout=-1,
            stderr=-1,
            env=None
        )

        self.assertEqual(compiled_file_path, out_dir / "main.json")