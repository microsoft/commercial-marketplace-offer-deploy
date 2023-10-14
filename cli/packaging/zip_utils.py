import os
from pathlib import Path
import zipfile
import shutil


# create a zip archive class called ZipArchive

def unzip_file(file_name, out_dir):
  path = Path(file_name).resolve()
  with zipfile.ZipFile(path, 'r') as zip_file:
    zip_file.extractall(out_dir)


def zip_dir(dir_path, file_path):
  file = Path(file_path)
  archive_path = str(file.resolve()).replace(file.suffix, '')

  shutil.make_archive(archive_path, 'zip', Path(dir_path), '')
  
  rename_file = Path(archive_path + '.zip')
  rename_file.rename(rename_file.with_suffix(file.suffix))

  return rename_file.with_suffix(file.suffix)