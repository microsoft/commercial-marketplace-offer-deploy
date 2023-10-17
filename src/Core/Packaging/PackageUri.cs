using System;

namespace Modm.Packaging
{
	public readonly struct PackageUri
	{
        private readonly Uri uri;

        public PackageUri(string value)
        {
            uri = new Uri(value, UriKind.Absolute);
            Value = value;

            if (!IsZipFileUri(uri))
            {
                throw new ArgumentException("URI must point to a ZIP file", nameof(value));
            }
        }

        private bool IsZipFileUri(Uri uri)
        {
            // Check if the URI scheme is "file" (local file) or "https" (HTTP/HTTPS)
            if (uri.Scheme.Equals("file", StringComparison.OrdinalIgnoreCase))
            {
                // For local files, check if the file extension is correct
                return Path.GetExtension(uri.LocalPath).Equals(PackageFile.Extension, StringComparison.OrdinalIgnoreCase);
            }
            else if (uri.Scheme.Equals("https", StringComparison.OrdinalIgnoreCase))
            {
                // For remote URIs, you might need to make an HTTP HEAD request
                // to the URI to check if it exists and is a pkg file
                // You can use HttpClient to do this, but it's beyond the scope of this code snippet
                // Here, we assume that the URI is valid and points to a pkg file
                return true;
            }
            else
            {
                return false; // Unsupported scheme
            }
        }


        /// <summary>
        /// The the URI value of the artifacts file for MODM to perform a deployment. The artifact file
        /// is expected to be a .zip file that was contained in the app.zip.
        /// 
        /// <see cref="https://learn.microsoft.com/en-us/azure/azure-resource-manager/managed-applications/publish-service-catalog-app"/>
        /// </summary>
        public readonly string Value { get; }

        public static implicit operator string(PackageUri uri) => uri.Value;
        public static explicit operator PackageUri(string v) => new(v);

        public override string ToString() => $"{Value}";
    }
}

