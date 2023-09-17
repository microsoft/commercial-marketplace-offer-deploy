using System;
namespace Modm.Deployments
{
	public class DeploymentDefinition
	{
		public required string DeploymentType { get; set; }

        public string OfferName { get; set; }
        public string OfferDescription { get; set; }

        /// <summary>
        /// The relative path to the main template, e.g. templates/mainTemplate.tf, templates/mainTemplate.
        /// </summary>
        public required string MainTemplate { get; set; }
    }
}

