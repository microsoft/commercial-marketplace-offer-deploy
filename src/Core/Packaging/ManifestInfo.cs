using System;
namespace Modm.Packaging
{
    /// <summary>
    /// The content of a manifest file
    /// </summary>
	public class ManifestInfo
	{
        /// <summary>
        /// The relative path to the main template, e.g. template/main.tf
        /// </summary>
        public required string MainTemplate { get; set; }
        public required string DeploymentType { get; set; }
        public OfferInfo Offer { get; set; }
    }

    public record OfferInfo
    {
        public string Name { get; set; }
        public string Description { get; set; }
    }
}

