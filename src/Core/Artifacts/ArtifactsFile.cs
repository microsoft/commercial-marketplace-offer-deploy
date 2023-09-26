using System;
using System.IO;
using System.IO.Compression;
using System.Security.Cryptography;

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
            ChangeDirectoryPermissions(ExtractedTo);

            IsExtracted = true;
        }

        /// <summary>
        /// Validates that the file to a given hash string
        /// </summary>
        /// <param name="signature">The input string of the expected hash</param>
        /// <returns>Whether the hashes are the same</returns>
        public bool IsValidSignature(string signature)
        {
            var computedHash = ComputeSha256Hash(this.filePath);
            return computedHash.Equals(signature);
        }

        /// <summary>
        /// Computes the Sha256 Hash for the Artifacts file path
        /// </summary>
        /// <param name="filePath">The location of the ArtifactsUri</param>
        /// <returns></returns>
        private string ComputeSha256Hash(string filePath)
        {
            using (FileStream stream = File.OpenRead(filePath))
            {
                using (SHA256 sha256 = SHA256.Create())
                {
                    byte[] hashBytes = sha256.ComputeHash(stream);
                    return BitConverter.ToString(hashBytes).Replace("-", "").ToLower();
                }
            }
        }

        /// <summary>
        /// This method will change the permissions of the directory to allow everyone to have full control
        /// </summary>
        /// <param name="directoryPath"></param>
        /// <exception cref="DirectoryNotFoundException"></exception>
        public static void ChangeDirectoryPermissions(string directoryPath)
        {
            try
            {
                DirectoryInfo directoryInfo = new(directoryPath);

                var ownerPermissions = UnixFileMode.UserExecute | UnixFileMode.UserRead | UnixFileMode.UserWrite;
                var groupPermissions = UnixFileMode.GroupExecute | UnixFileMode.GroupRead | UnixFileMode.GroupWrite;
                var othersPermissions = UnixFileMode.OtherExecute | UnixFileMode.OtherRead | UnixFileMode.OtherWrite;
                var allPermissions = ownerPermissions | groupPermissions | othersPermissions;

                directoryInfo.UnixFileMode = allPermissions;
            }
            catch
            {

            }
        }
    }
}

