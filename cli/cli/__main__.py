
import json
import click
from .commands import build_application_package, build_function_app_package


@click.group()
def cli():
    """
    Commercial Marketplace Offer Deployment Manager (MODM) CLI that exposes application packaging
    Development Version: 0.1.0
    """
    pass

@cli.group('functionapp')
def functionapp():
    """
    Package commands
    """
    pass

functionapp.add_command(build_function_app_package)

@cli.group('package')
def package():
    """
    Package commands
    """
    pass

package.add_command(build_application_package)

if __name__ == '__main__':
    cli()