import unittest
from packaging import ApplicationPackage, DeploymentType

class TestApplicationPackage(unittest.TestCase):

  def test_main_template(self):
    app_package = ApplicationPackage("main.tf" , "fake_create_ui_definition")
    self.assertEqual(app_package.manifest.main_template, "main.tf")
    self.assertEqual(app_package.manifest.deployment_type, DeploymentType.terraform)

  def test_get_main_template(self):
    app_package = ApplicationPackage("", "")
    self.assertEqual(app_package.main_template.name, "mainTemplate.json")
    self.assertIsNotNone(app_package.main_template.document)
    self.assertIsNotNone(app_package.main_template.document['variables']['userDataObject'])

  # def test_create_ui_definition(self):
  #   app_package = ApplicationPackage("deployment_type", "main_template")
  #   app_package.create_ui_definition("new_ui_definition")
  #   self.assertEqual(app_package.create_ui_definition, "new_ui_definition")

  # def test_create(self):
  #   app_package = ApplicationPackage("deployment_type", "main_template")
  #   app_package.create()
  #   # Add assertions to check that the package was created successfully
