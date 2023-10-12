from zipfile import ZipFile


class ApplicationPackageBuilder:
  def build(self):
    print('Building application package...')
    with ZipFile('spam.zip', 'w') as myzip:
      myzip.write('eggs.txt')

  def _create_installer_zip(self):
    with ZipFile('installer.zip', 'w') as myzip:
      myzip.write('eggs.txt')