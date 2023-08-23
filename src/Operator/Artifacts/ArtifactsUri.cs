using System;
namespace Operator.Artifacts
{
	public record ArtifactsUri
	{
		public const string ContentPath = "content/";

        /// <summary>
        /// The base URI to the blob storage container that is holding all the artifacts
        /// from the app.zip file
        /// <see cref="https://learn.microsoft.com/en-us/azure/azure-resource-manager/managed-applications/publish-service-catalog-app"/>
        /// </summary>
        public required string Container { get; set; }

		/// <summary>
		/// The URI
		/// </summary>
		public required string Content { get; set; }
	}
}

