namespace Modm
{
	class ControllerOptions
	{
		public string? MachineName { get; set; }
		public string? Fqdn { get; set; }
		public string? ComposeFilePath { get; set; }

		public string ComposeFileDirectory
		{
			get
			{
				if (ComposeFilePath == null)
				{
					throw new InvalidOperationException("Docker compose file path cannot be null.");
				}

				var file = new FileInfo(ComposeFilePath);

				if (!file.Exists || file.DirectoryName == null)
				{
					throw new InvalidOperationException("Invalid file or directory path for docker compose file.");
				}
				return file.DirectoryName;

            }
		}
    }
}

