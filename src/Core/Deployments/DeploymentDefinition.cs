using System;
namespace Modm.Deployments
{
    /// <summary>
    /// Represents the contents of the manifest file contained in the artifacts archive file, e.g. the content.zip inside the app.zip
    /// </summary>
	public class DeploymentDefinition
	{
        /// <summary>
        /// The relative path to the main template, e.g. template/main.tf
        /// </summary>
        public required string MainTemplate { get; set; }
        public required string DeploymentType { get; set; }
        public OfferInfo Offer { get; set; }
    }

    public class OfferInfo
    {
        public string Name { get; set; }
        public string Description { get; set; }
    }
}

