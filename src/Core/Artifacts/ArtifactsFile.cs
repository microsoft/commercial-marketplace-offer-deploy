using System;
using System.IO.Compression;
using System.Security.AccessControl;
using System.Runtime.InteropServices;

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
                Directory.Delete(ExtractedTo);
            }

            // because unzip will use the name of the zip file when extracting
            // unzip directly to ./content
            var destinationDirectoryName = Path.GetDirectoryName(filePath);
            
            ZipFile.ExtractToDirectory(filePath, destinationDirectoryName, overwriteFiles: true);
            ChangeDirectoryPermissions(destinationDirectoryName);

            IsExtracted = true;
        }

        /// <summary>
        /// This method will change the permissions of the directory to allow everyone to have full control
        /// </summary>
        /// <param name="directoryPath"></param>
        /// <exception cref="DirectoryNotFoundException"></exception>
        public void ChangeDirectoryPermissions(string directoryPath)
        {
            DirectoryInfo directoryInfo = new DirectoryInfo(directoryPath);
            var ownerPermissions = UnixFileMode.UserExecute | UnixFileMode.UserRead | UnixFileMode.UserWrite;
            var groupPermissions = UnixFileMode.GroupExecute | UnixFileMode.GroupRead | UnixFileMode.GroupWrite;
            var othersPermissions = UnixFileMode.OtherExecute | UnixFileMode.OtherRead | UnixFileMode.OtherWrite;
            var allPermissions = ownerPermissions | groupPermissions | othersPermissions;
            directoryInfo.UnixFileMode = allPermissions;
        }

        /// <summary>
        /// Reads and returns the manifest file as the deployment definition
        /// </summary>
        /// <returns></returns>
        public async Task<DeploymentDefinition> ReadManifestFile()
        {
            if (!IsExtracted)
            {
                Extract();
            }

            var definition = await ManifestFile.Read(ExtractedTo);
            return definition;
        }
    }
}

