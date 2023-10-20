import json
from pathlib import Path
import shutil
import tempfile
import click
from packaging import ApplicationPackage, _zip_utils
from packaging.application_package import CreateApplicationPackageOptions
from packaging.azure import ArmTemplate
import os

from packaging.function_app_package import FunctionAppPackage

# pylint: disable=length-too-long
FUNCTION_APP_PACKAGE_URL = ''

@click.command('build')
@click.option("-v", "--vmi-reference-id", help="The VMI reference id", required=True)
@click.option("-t", "--template-file", help="The path to the application's main template.", required=True)
@click.option("-u", "--create-ui-definition", help="The path to the createUiDefinition.json file", required=True)
@click.option("-o", "--out-dir", help="The location where the application package will be created")
@click.option("-n", "--name", help="The name of the application")
@click.option("-d", "--description", help="The description of the application")
@click.argument('current_working_dir', type=click.Path(exists=True))
def build_application_package(vmi_reference_id, template_file, create_ui_definition, name, description, current_working_dir, out_dir = None):
    """Builds an application package and produces an app.zip"""
    cwd = Path(current_working_dir)
    out_dir = cwd.joinpath(out_dir).resolve() if out_dir else cwd.joinpath(".").resolve()
    
    resolved_template_file = cwd.joinpath(template_file).resolve()
    resolved_create_ui_definition = cwd.joinpath(create_ui_definition).resolve()

    package = ApplicationPackage(resolved_template_file, resolved_create_ui_definition, name, description)

    function_app_package = FunctionAppPackage.from_resource()
    options = CreateApplicationPackageOptions(vmi_reference_id, function_app_package)
    
    result = package.create(options, out_dir)

    click.echo(json.dumps(result.serialize(), indent=2))


@click.command('build')
@click.option("-f", "--csproj-file", help="Path to the .csproj file of the Function App", required=True)
@click.option("-o", "--out-dir", help="The location where the zip will be placed")
@click.argument('current_working_dir', type=click.Path(exists=True))
def build_function_app_package(csproj_file, current_working_dir, out_dir = None):
    """Builds a function.zip"""
    cwd = Path(current_working_dir)
    out_dir = _resolve_path(cwd, out_dir)
    csproj_file = _resolve_path(cwd, csproj_file)
    temp_dir = tempfile.mkdtemp()

    import subprocess
    process = subprocess.Popen(['dotnet', 'publish', csproj_file, '-c', 'Release', '-o', temp_dir], 
                           stdout=subprocess.PIPE,
                           universal_newlines=True)

    click.echo("Building function app...")

    while True:
        output = process.stdout.readline()
        click.echo(output.strip())
        # Do something else
        return_code = process.poll()
        if return_code is not None:
            # Process has finished, read rest of the output 
            for output in process.stdout.readlines():
                click.echo(output.strip())
            break
    
    click.echo("Creating function app package.")
    _zip_utils.zip_dir(temp_dir, out_dir / 'function.zip')

    shutil.rmtree(temp_dir)
    click.echo("Package created.")
    click.echo(out_dir / 'function.zip')


def _resolve_path(current_working_dir, path):
    cwd = Path(current_working_dir)
    return cwd.joinpath(path).resolve() if path else cwd.joinpath(".").resolve()
    