import hashlib
import json
import os
from pathlib import Path
import shutil
import tarfile
import tempfile
import click
from modm import _zip_utils
from modm.marketplace.application_packaging_options import ApplicationPackageOptions
from modm.marketplace.application_package_info import ApplicationPackageInfo
from modm.marketplace.application_package import new_application_package
from modm.installer.client_app_package import ClientAppPackage
from modm.installer.version import InstallerVersion


@click.help_option("-h", "--help")
@click.command("build")
@click.option("-n", "--name", help="The name of the application")
@click.option("-d", "--description", help="The description of the application")
@click.option("-v", "--version", default="latest", help="The version of the installer (modm) to package the application with.")
@click.option("--vmi-reference", default=False, type=bool, help="Whether to reference the VMI directly when packaging.")
@click.option("--vmi-reference-id", default=None, type=str, help="The ID of the VMI reference to use to override the published reference.")
@click.option(
    "--resources-file",
    default=None,
    help="The reference to resources file to use when package. This must be used in combination with a direct vmi-reference-id being set.",
)
@click.option("-f", "--main-template", help="The path to the application's main template.", required=True)
@click.option("--create-ui-definition", help="The path to the createUiDefinition.json file", required=True)
@click.option("-o", "--out-dir", help="The location where the application package will be created", required=True)
@click.argument("current_working_dir", type=click.Path(exists=True))
def build_application_package(
    name,
    description,
    version,
    vmi_reference,
    vmi_reference_id,
    resources_file,
    main_template,
    create_ui_definition,
    current_working_dir,
    out_dir=None,
):
    """Builds an application package and produces an app.zip"""
    cwd = Path(current_working_dir)
    out_dir = _resolve_path(cwd, out_dir)

    resolved_template_file = cwd.joinpath(main_template).resolve()
    resolved_create_ui_definition = cwd.joinpath(create_ui_definition).resolve()

    if resources_file is not None:
        resources_file = cwd.joinpath(resources_file).resolve()

    info = ApplicationPackageInfo(resolved_template_file, resolved_create_ui_definition, name, description)
    options = ApplicationPackageOptions(version, vmi_reference, vmi_reference_id, resources_file, out_dir)

    package = new_application_package()
    result = package.create(info, options)

    click.echo(json.dumps(result.serialize(), indent=2))


@click.help_option("-h", "--help")
@click.command("create-function-app-package")
@click.option("-f", "--csproj-file", help="Path to the .csproj file of the client app", required=True)
@click.option("-o", "--out-dir", help="The location where the zip will be placed")
@click.argument("current_working_dir", type=click.Path(exists=True))
def create_client_app_package(csproj_file, current_working_dir, out_dir=None):
    """creates a clientapp.zip"""
    _create_client_app_package(csproj_file, current_working_dir, out_dir)


@click.help_option("-h", "--help")
@click.command("create-resources-tarball")
@click.option("-v", "--version", default=None, type=str, help="The version of the installer (modm) to package the application with.")
@click.option(
    "-t", "--templates-dir", help="The location where the templates file is located (mainTemplate.json and viewDefinition.json)", required=True
)
@click.option("-f", "--csproj-file", help="Path to the .csproj file of the client app", required=True)
@click.option("-o", "--out-dir", help="The location where the tarball will be created", required=True)
@click.argument("current_working_dir", type=click.Path(exists=True))
def create_resources_tarball(version, templates_dir, csproj_file, current_working_dir, out_dir=None):
    cwd = Path(current_working_dir)
    out_dir = _resolve_path(cwd, out_dir)

    main_template_file = _resolve_path(cwd, templates_dir) / "mainTemplate.json"
    view_definition_file = _resolve_path(cwd, templates_dir) / "viewDefinition.json"
    create_ui_definition_file = _resolve_path(cwd, templates_dir) / "createUiDefinition.json"
    client_app_file = _create_client_app_package(csproj_file, current_working_dir, out_dir)

    if version is not None:
        installer_version = InstallerVersion(version)
        out_file = out_dir / f"resources-{installer_version.name}.tar.gz"
    else:
        out_file = out_dir / f"resources.tar.gz"

    click.echo(f"Creating tarball {out_file}...")

    with tarfile.open(out_file, "w:gz") as tar:
        tar.add(main_template_file, arcname=main_template_file.name)
        tar.add(view_definition_file, arcname=view_definition_file.name)
        tar.add(create_ui_definition_file, arcname=create_ui_definition_file.name)
        tar.add(client_app_file, arcname=client_app_file.name)

    click.echo(f"resources '{out_file.name}' created.")

    if version is not None:
        result = {}
        result[
            "downloadUrl"
        ] = f"https://github.com/microsoft/commercial-marketplace-offer-deploy/releases/download/{installer_version.name}/{out_file.name}"
        result["filename"] = out_file.name
        result["sha256Digest"] = _get_sha256_digest(out_file)
    else:
        result = {}
        result["filename"] = str(out_file)
        result["sha256Digest"] = _get_sha256_digest(out_file)

    click.echo(json.dumps(result, indent=2))


def _create_client_app_package(csproj_file, current_working_dir, out_dir=None):
    """creates a clientapp.zip"""
    cwd = Path(current_working_dir)
    out_dir = _resolve_path(cwd, out_dir)
    csproj_file = _resolve_path(cwd, csproj_file)
    temp_dir = tempfile.mkdtemp()

    import subprocess

    process = subprocess.Popen(["dotnet", "publish", csproj_file, "-c", "Release", "-o", temp_dir], stdout=subprocess.PIPE, universal_newlines=True)

    click.echo("Creating client app package...")

    while True:
        output = process.stdout.readline()
        click.echo("  " + output.strip())
        # Do something else
        return_code = process.poll()
        if return_code is not None:
            # Process has finished, read rest of the output
            for output in process.stdout.readlines():
                click.echo("  " + output.strip())
            break

    click.echo("Creating client app package.")

    out_file = out_dir / ClientAppPackage.file_name
    _zip_utils.zip_dir(temp_dir, out_file)

    shutil.rmtree(temp_dir)
    click.echo(out_file)

    return out_file


def _resolve_path(current_working_dir, path):
    cwd = Path(current_working_dir)
    return cwd.joinpath(path).resolve() if path else cwd.joinpath(".").resolve()


def _get_sha256_digest(file):
    hash_sha256 = hashlib.sha256()
    with open(file, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            hash_sha256.update(chunk)
    return hash_sha256.hexdigest()
