using System;

namespace Modm.Artifacts
{
	public readonly struct ArtifactsUri
	{
        public const string ArtifactsFileExtension = "zip";

        private readonly Uri uri;

        public ArtifactsUri(string value)
        {
            uri = new Uri(value, UriKind.Absolute);
            Value = value;

            if (!uri.IsFile)
            {
                throw new ArgumentException("must be a URI to a zip file", nameof(value));
            }
        }

        /// <summary>
        /// The the URI value of the artifacts file for MODM to perform a deployment. The artifact file
        /// is expected to be a .zip file that was contained in the app.zip.
        /// 
        /// <see cref="https://learn.microsoft.com/en-us/azure/azure-resource-manager/managed-applications/publish-service-catalog-app"/>
        /// </summary>
        public readonly string Value { get; }

        public readonly string FileName
        {
            get
            {
                return Path.GetFileName(uri.LocalPath);
            }
        }

        public static implicit operator string(ArtifactsUri uri) => uri.Value;
        public static explicit operator ArtifactsUri(string v) => new(v);

        public override string ToString() => $"{Value}";
    }
}

