using System;
namespace Modm.Artifacts
{
	public class ArtifactsDownloadOptions
	{
        public const string ConfigSectionKey = "ArtifactsDownload";

        public required string SavePath { get; set; }
    }
}

