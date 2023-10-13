import click

@click.command()
@click.option("-t", "--main-template", default="", help="The path to the application's main template.")
@click.option("-u", "--create-ui-definition", prompt="Your name", help="The path to the createUiDefinition.json file")
def build(count, name):
    """Builds an application package, e.g. an app.zip"""
    for _ in range(count):
        click.echo(f"Hello, {name}!")