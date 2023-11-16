import re


class Version:
    def __init__(self, *args):
        # document the constructor options in the docstring
        """
        Constructor for Version class. The constructor can be called in the following ways:

        1. Version("v0.0.0") or Version("v0.0.0-suffix") or Version("0.0.0")
        2. Version(0, 0, 0)
        3. Version(0, 0, 0, "suffix")

        The constructor will also parse the version name and set the major, minor, patch, and suffix properties
        accordingly. The version name option must be in the following format: v0.0.0 or v0.0.0-<suffix>
        """

        # am empty version
        self._major = 0
        self._minor = 0
        self._patch = 0
        self._suffix = ""

        if len(args) == 1 and isinstance(args[0], str):
            self.major, self.minor, self.patch, self.suffix = self.parse(args[0])
        elif len(args) == 1 and isinstance(args[0], Version):
            self.major = args[0].major
            self.minor = args[0].minor
            self.patch = args[0].patch
            self.suffix = args[0].suffix
        elif len(args) == 3:
            self.major = args[0]
            self.minor = args[1]
            self.patch = args[2]
            self.suffix = ""
        elif len(args) == 4:
            self.major = args[0]
            self.minor = args[1]
            self.patch = args[2]
            self.suffix = args[3]

    @property
    def major(self):
        return self._major

    @major.setter
    def major(self, value):
        self._major = int(value)

    @property
    def minor(self):
        return self._minor

    @minor.setter
    def minor(self, value):
        self._minor = int(value)

    @property
    def patch(self):
        return self._patch

    @patch.setter
    def patch(self, value):
        self._patch = int(value)

    @property
    def suffix(self):
        return self._suffix

    @suffix.setter
    def suffix(self, value):
        self._suffix = value

    @property
    def name(self):
        name = f"v{self.major}.{self.minor}.{self.patch}"
        if self.suffix != "":
            name += f"-{self.suffix}"
        return name

    def parse(self, name: str):
        if name == "" or name is None:
            raise ValueError(f"Invalid version name format: '{name}'. Expected format: v0.0.0 or v0.0.0-<suffix>")

        version_parts = re.sub("-[a-zA-Z0-9]*", "", name).removeprefix("v").split(".")
        if len(version_parts) != 3:
            raise ValueError(f"Invalid version name format: '{name}'. Expected major.minor.patch format")

        major = version_parts[0]
        minor = version_parts[1]
        build = version_parts[2]
        suffix = name.split("-")[1] if len(name.split("-")) == 2 else ""

        return major, minor, build, suffix

    def __eq__(self, __value: object) -> bool:
        if isinstance(__value, str):
            value = Version(__value)
            return self.major == value.major and self.minor == value.minor and self.patch == value.patch and self.suffix == value.suffix
        if isinstance(__value, Version):
            return self.major == __value.major and self.minor == __value.minor and self.patch == __value.patch and self.suffix == __value.suffix
        else:
            return False