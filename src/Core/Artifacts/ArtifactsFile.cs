using System.IO.Compression;
using Modm.Deployments;

namespace Modm.Artifacts
{
    /// <summary>
    /// The content.zip file that represents the artifacts used for the installation by MODM
    /// </summary>
	public class ArtifactsFile
    {
        /// <summary>
        /// the name of the content file that must always be used.
        /// This file will be included in the app.zip file at the root
        /// </summary>
        public const string FileName = "content.zip";

        /// <summary>
        /// the default directory name where the content file will be extracted to
        /// </summary>
        public const string DestinationDirectoryName = "content";

        private readonly string filePath;

        /// <summary>
        /// Gets the extracted to directory path
        /// </summary>
        public string ExtractedTo => Path.Combine(new FileInfo(filePath).DirectoryName, DestinationDirectoryName);

        public bool IsExtracted { get; private set; }

        public ArtifactsFile(string filePath)
        {
            this.filePath = filePath;
            IsExtracted = false;
        }

        /// <summary>
        /// Extracts the contents of the content.zip to a directory <see cref="DestinationDirectoryName"/> next to it
        /// and returns the directory path
        /// </summary>
        public void Extract()
		{
            if (IsExtracted)
                return;

            if (Directory.Exists(ExtractedTo))
            {
                Directory.Delete(ExtractedTo, recursive: true);
            }

            // because unzip will use the name of the zip file when extracting
            // unzip directly to ./content
            var destinationDirectoryName = Path.GetDirectoryName(filePath);
            ZipFile.ExtractToDirectory(filePath, destinationDirectoryName, overwriteFiles: true);

            IsExtracted = true;
        }
    }
}

