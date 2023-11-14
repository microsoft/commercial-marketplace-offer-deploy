from msrest.serialization import Model


class ApplicationPackageResult(Model):
    _attribute_map = {
        "file": {"key": "file", "type": "str"},
        "validation_results": {"key": "validationResults", "type": "[object]"},
    }

    def __init__(self, **kwargs) -> None:
        self.file = None
        self.validation_results = kwargs.get("validation_results", [])
        self._installer_package = kwargs.get("installer_package", None)
        self._client_app_name = kwargs.get("client_app_name", None)

    @property
    def client_app_name(self):
        return self._client_app_name

    @property
    def installer_package(self):
        return self._installer_package