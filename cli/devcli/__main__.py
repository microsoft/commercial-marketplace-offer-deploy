
import json
import click
from .commands import build_application_package, create_client_app_package, create_resources_archive


@click.group()
def cli():
    """
    Commercial Marketplace Offer Deployment Manager (MODM) CLI that exposes application packaging
    Development Version: 0.1.0
    """
    pass

@cli.group('util')
def util():
    """
    Utility commands
    """
    pass

util.add_command(create_client_app_package)
util.add_command(create_resources_archive)

@cli.group('package')
def package():
    """
    Package commands
    """
    pass

package.add_command(build_application_package)

if __name__ == '__main__':
    cli()