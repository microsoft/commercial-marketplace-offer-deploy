// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿namespace Modm
{
	class ControllerOptions
	{
		public string? MachineName { get; set; }
		public string? Fqdn { get; set; }
		public string? ComposeFilePath { get; set; }
		public string? StateFilePath { get; set; }

		public string ComposeFileDirectory
		{
			get
			{
				return GetFileDirectory(ComposeFilePath);
            }
		}

		public string StateFileDirectory
		{
			get
			{
				return GetFileDirectory(StateFilePath);
			}
		}

		private string GetFileDirectory(string? filePath)
		{
            if (string.IsNullOrEmpty(filePath))
            {
                throw new InvalidOperationException($"Filepath for path cannot be null or empty.");
            }

            var file = new FileInfo(filePath);

            if (!file.Exists || file.DirectoryName == null)
            {
                throw new InvalidOperationException("Invalid file or directory path for docker compose file.");
            }
            return file.DirectoryName;

        }
    }
}

