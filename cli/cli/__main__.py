
import json
import click
from .commands import build


@click.group()
def cli():
    """
    Commercial Marketplace Offer Deployment Manager (MODM) CLI that exposes application packaging
    Development Version: 0.1.0
    """
    pass

cli.add_command(build)

if __name__ == '__main__':
    cli()