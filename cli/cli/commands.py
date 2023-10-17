import json
from pathlib import Path
import click
from packaging import ApplicationPackage
from packaging.application_package import CreateApplicationPackageOptions
from packaging.azure import ArmTemplate
import os

from packaging.function_app_package import FunctionAppPackage

# pylint: disable=length-too-long
FUNCTION_APP_PACKAGE_URL = ''

@click.command()
@click.option("-v", "--vmi-reference-id", help="The VMI reference id", required=True)
@click.option("-t", "--template-file", help="The path to the application's main template.", required=True)
@click.option("-u", "--create-ui-definition", help="The path to the createUiDefinition.json file", required=True)
@click.option("-o", "--out-dir", help="The location where the application package will be created")
@click.option("-n", "--name", help="The name of the application")
@click.option("-d", "--description", help="The description of the application")
@click.argument('current_working_dir', type=click.Path(exists=True))
def build(vmi_reference_id, template_file, create_ui_definition, name, description, current_working_dir, out_dir = None):
    """Builds an application package and produces an app.zip"""
    cwd = Path(current_working_dir)
    out_dir = cwd.joinpath(out_dir).resolve() if out_dir else cwd.joinpath(".").resolve()
    
    # module will get executed under ./cli dir context, so let's go up one level
    resolved_template_file = cwd.joinpath(template_file).resolve()
    resolved_create_ui_definition = cwd.joinpath(create_ui_definition).resolve()

    package = ApplicationPackage(resolved_template_file, resolved_create_ui_definition, name, description)

    function_app_package = FunctionAppPackage.from_uri(FUNCTION_APP_PACKAGE_URL, out_dir)

    options = CreateApplicationPackageOptions(vmi_reference_id, function_app_package)
    result = package.create(options, out_dir)

    click.echo(json.dumps(result.serialize(), indent=2))
    