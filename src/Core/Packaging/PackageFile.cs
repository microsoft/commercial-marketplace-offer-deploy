using System;
using System.IO;
using System.IO.Compression;
using System.Security.Cryptography;
using Microsoft.Extensions.Logging;
using Modm.IO;

namespace Modm.Packaging
{
    /// <summary>
    /// The installer.pkg file that represents the artifacts used for the installation of an app by modm
    /// </summary>
	public class PackageFile
    {
        public const string Extension = ".zip";

        public const string HashAttributeName = "packageHash";

        /// <summary>
        /// the name of the installer package file that must always be used.
        /// This file will be included in the app.zip file at the root
        /// </summary>
        public const string FileName = "installer.zip";

        /// <summary>
        /// the default directory name where the installer package file will be extracted to
        /// </summary>
        public const string DestinationDirectoryName = "installer";

        private readonly string filePath;
        private readonly ILogger<PackageFile> logger;

        /// <summary>
        /// Gets the extracted to directory path
        /// </summary>
        public string ExtractedTo => Path.Combine(new FileInfo(filePath).DirectoryName, DestinationDirectoryName);

        public bool IsExtracted { get; private set; }

        public string ComputedHash
        {
            get
            {
                return ComputeSha256Hash(this.filePath);
            }
        }

        internal PackageFile(string filePath, ILogger<PackageFile> logger)
        {
            this.filePath = filePath;
            this.logger = logger;
            IsExtracted = false;
        }

        /// <summary>
        /// Extracts the contents of the installer package to a directory <see cref="DestinationDirectoryName"/> next to it
        /// and returns the directory path
        /// </summary>
        public void Extract()
		{
            if (IsExtracted)
                return;

            if (Directory.Exists(ExtractedTo))
            {
                logger.LogDebug("directory existed prior to extraction. Deleting");
                Directory.Delete(ExtractedTo, recursive: true);
            }

            // because unzip will use the name of the zip file when extracting, unzip directly to ./installer
            var destinationDirectoryName = Path.GetDirectoryName(filePath);

            ZipFile.ExtractToDirectory(filePath, destinationDirectoryName, overwriteFiles: true);
            DirectoryPermissions.AllowFullControl(ExtractedTo);

            IsExtracted = true;
        }

        /// <summary>
        /// Validates that the file to a given hash string
        /// </summary>
        /// <param name="hash">The input string of the expected hash</param>
        /// <returns>Whether the hashes are the same</returns>
        public bool IsValidHash(string hash)
        {
            if (hash == null)
            {
                return false;
            }

            var computedHash = ComputeSha256Hash(this.filePath);
            return computedHash.Equals(hash);
        }

        /// <summary>
        /// Computes the Sha256 Hash for the Artifacts file path
        /// </summary>
        /// <param name="filePath">The location of the ArtifactsUri</param>
        /// <returns></returns>
        private static string ComputeSha256Hash(string filePath)
        {
            using FileStream stream = File.OpenRead(filePath);
            using SHA256 sha256 = SHA256.Create();

            byte[] hashBytes = sha256.ComputeHash(stream);
            return BitConverter.ToString(hashBytes).Replace("-", "").ToLower();
        }
    }
}

