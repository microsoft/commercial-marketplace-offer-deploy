from pathlib import Path
import requests
import os


def download_file(file_url: str, file_path: Path) -> str:
    """
    Downloads a file from the given URL and saves it to the given file path.
    Returns the file path if successful, raises an exception otherwise.
    """
    try:
        response = requests.get(file_url, stream=True)
        response.raise_for_status()

        # Create the directory if it doesn't exist but remove file so we overwrite it
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        if file_path.exists():
            os.remove(file_path)

        with open(file_path, "wb") as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)
        return file_path

    except Exception as e:
        raise Exception(f"Failed to download file from {file_url}: {e}")
