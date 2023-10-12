import click

@click.command()
@click.option("-t", "--main-template", default="", help="The application's main template.")
@click.option("--name", prompt="Your name", help="The person to greet.")
def build(count, name):
    """Builds an application package, e.g. an app.zip"""
    for _ in range(count):
        click.echo(f"Hello, {name}!")