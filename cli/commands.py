import click
from cli.packaging import ApplicationPackage

@click.command()
@click.option("-t", "--template", prompt="Application template file path", help="The path to the application's main template.")
@click.option("-u", "--create-ui-definition", prompt="createUiDefinition.json path", help="The path to the createUiDefinition.json file")
@click.option("-n", "--name", prompt="Application name", help="The name of the application")
@click.option("-d", "--description", prompt="Application description", help="The description of the application")
def build(template, create_ui_definition, name, description):
    """Builds an application package, e.g. an app.zip"""
    package = ApplicationPackage(template, create_ui_definition, name, description)
    result = package.create()

    print(result.serialize())
    